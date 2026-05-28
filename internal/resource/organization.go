package resource

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/simfra-dev/terraform-provider-simfra/internal/awsclient"
	"github.com/simfra-dev/terraform-provider-simfra/internal/providerdata"
)

var _ resource.Resource = &organizationResource{}

type organizationResource struct {
	aws *awsclient.Client
}

type organizationModel struct {
	ID               types.String `tfsdk:"id"`
	RootID           types.String `tfsdk:"root_id"`
	FeatureSet       types.String `tfsdk:"feature_set"`
	ARN              types.String `tfsdk:"arn"`
	MasterAccountID  types.String `tfsdk:"master_account_id"`
	MasterAccountARN types.String `tfsdk:"master_account_arn"`
}

func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

func (r *organizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an AWS Organization with optional specific org and root IDs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Organization ID (e.g. o-abc1234567). If specified, the organization is created with this exact ID. If omitted, Simfra generates one.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"root_id": schema.StringAttribute{
				Description: "Organization root ID (e.g. r-ab12). If specified, the root is created with this exact ID. If omitted, Simfra generates one.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"feature_set": schema.StringAttribute{
				Description: "Feature set: ALL or CONSOLIDATED_BILLING.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ALL"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"arn": schema.StringAttribute{
				Description: "Organization ARN.",
				Computed:    true,
			},
			"master_account_id": schema.StringAttribute{
				Description: "Management account ID.",
				Computed:    true,
			},
			"master_account_arn": schema.StringAttribute{
				Description: "Management account ARN.",
				Computed:    true,
			},
		},
	}
}

func (r *organizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerdata.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *providerdata.ProviderData, got: %T", req.ProviderData))
		return
	}
	if data.AWS == nil {
		resp.Diagnostics.AddError("AWS credentials required", "simfra_organization requires access_key and secret_key in the provider configuration.")
		return
	}
	r.aws = data.AWS
}

func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &organizations.CreateOrganizationInput{
		FeatureSet: orgtypes.OrganizationFeatureSet(plan.FeatureSet.ValueString()),
	}

	overrides := make(map[string]string)
	if !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		overrides["OrganizationId"] = plan.ID.ValueString()
	}
	if !plan.RootID.IsNull() && !plan.RootID.IsUnknown() {
		overrides["RootId"] = plan.RootID.ValueString()
	}
	if len(overrides) > 0 {
		ctx = awsclient.WithOverrides(ctx, overrides)
	}

	out, err := r.aws.Organizations.CreateOrganization(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create organization", err.Error())
		return
	}

	plan.ID = types.StringValue(aws.ToString(out.Organization.Id))
	plan.ARN = types.StringValue(aws.ToString(out.Organization.Arn))
	plan.MasterAccountID = types.StringValue(aws.ToString(out.Organization.MasterAccountId))
	plan.MasterAccountARN = types.StringValue(aws.ToString(out.Organization.MasterAccountArn))
	plan.FeatureSet = types.StringValue(string(out.Organization.FeatureSet))

	// Get root ID from ListRoots
	roots, err := r.aws.Organizations.ListRoots(ctx, &organizations.ListRootsInput{})
	if err == nil && len(roots.Roots) > 0 {
		plan.RootID = types.StringValue(aws.ToString(roots.Roots[0].Id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.aws.Organizations.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		if isAWSNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read organization", err.Error())
		return
	}

	state.ID = types.StringValue(aws.ToString(out.Organization.Id))
	state.ARN = types.StringValue(aws.ToString(out.Organization.Arn))
	state.MasterAccountID = types.StringValue(aws.ToString(out.Organization.MasterAccountId))
	state.MasterAccountARN = types.StringValue(aws.ToString(out.Organization.MasterAccountArn))
	state.FeatureSet = types.StringValue(string(out.Organization.FeatureSet))

	roots, err := r.aws.Organizations.ListRoots(ctx, &organizations.ListRootsInput{})
	if err == nil && len(roots.Roots) > 0 {
		state.RootID = types.StringValue(aws.ToString(roots.Roots[0].Id))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	_, err := r.aws.Organizations.DeleteOrganization(ctx, &organizations.DeleteOrganizationInput{})
	if err != nil && !isAWSNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete organization", err.Error())
	}
}
