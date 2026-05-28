package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/simfra-dev/terraform-provider-simfra/internal/awsclient"
	"github.com/simfra-dev/terraform-provider-simfra/internal/providerdata"
)

var _ resource.Resource = &route53ZoneResource{}

type route53ZoneResource struct {
	aws *awsclient.Client
}

type route53ZoneModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Comment     types.String `tfsdk:"comment"`
	NameServers types.List   `tfsdk:"name_servers"`
}

func NewRoute53ZoneResource() resource.Resource {
	return &route53ZoneResource{}
}

func (r *route53ZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route53_zone"
}

func (r *route53ZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Route53 hosted zone with an optional specific zone ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Hosted zone ID. If specified, the zone is created with this exact ID. If omitted, Simfra generates one.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Domain name for the hosted zone (e.g. example.com).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the hosted zone.",
				Optional:    true,
			},
			"name_servers": schema.ListAttribute{
				Description: "Delegation name servers for the hosted zone.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *route53ZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerdata.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *providerdata.ProviderData, got: %T", req.ProviderData))
		return
	}
	if data.AWS == nil {
		resp.Diagnostics.AddError("AWS credentials required", "simfra_route53_zone requires access_key and secret_key in the provider configuration.")
		return
	}
	r.aws = data.AWS
}

func (r *route53ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan route53ZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &route53.CreateHostedZoneInput{
		Name:            aws.String(plan.Name.ValueString()),
		CallerReference: aws.String(uuid.New().String()),
	}
	if !plan.Comment.IsNull() && plan.Comment.ValueString() != "" {
		input.HostedZoneConfig = &route53types.HostedZoneConfig{
			Comment: aws.String(plan.Comment.ValueString()),
		}
	}

	if !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		ctx = awsclient.WithOverrides(ctx, map[string]string{"Id": plan.ID.ValueString()})
	}

	out, err := r.aws.Route53.CreateHostedZone(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create hosted zone", err.Error())
		return
	}

	zoneID := aws.ToString(out.HostedZone.Id)
	zoneID = strings.TrimPrefix(zoneID, "/hostedzone/")
	plan.ID = types.StringValue(zoneID)

	if out.DelegationSet != nil {
		ns, diags := types.ListValueFrom(ctx, types.StringType, out.DelegationSet.NameServers)
		resp.Diagnostics.Append(diags...)
		plan.NameServers = ns
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *route53ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state route53ZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.aws.Route53.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: aws.String(state.ID.ValueString()),
	})
	if err != nil {
		if isAWSNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read hosted zone", err.Error())
		return
	}

	state.Name = types.StringValue(strings.TrimSuffix(aws.ToString(out.HostedZone.Name), "."))
	if out.HostedZone.Config != nil && out.HostedZone.Config.Comment != nil {
		state.Comment = types.StringValue(aws.ToString(out.HostedZone.Config.Comment))
	}
	if out.DelegationSet != nil {
		ns, diags := types.ListValueFrom(ctx, types.StringType, out.DelegationSet.NameServers)
		resp.Diagnostics.Append(diags...)
		state.NameServers = ns
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *route53ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan route53ZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *route53ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state route53ZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.aws.Route53.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
		Id: aws.String(state.ID.ValueString()),
	})
	if err != nil && !isAWSNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete hosted zone", err.Error())
	}
}
