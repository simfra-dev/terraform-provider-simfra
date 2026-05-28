package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/simfra-dev/terraform-provider-simfra/internal/client"
	"github.com/simfra-dev/terraform-provider-simfra/internal/providerdata"
)

var _ datasource.DataSource = &portForwardTargetsDataSource{}

type portForwardTargetsDataSource struct {
	client *client.Client
}

type portForwardTargetsDataSourceModel struct {
	Targets []portForwardTargetModel `tfsdk:"targets"`
}

type portForwardTargetModel struct {
	ARN         types.String `tfsdk:"arn"`
	Service     types.String `tfsdk:"service"`
	ResourceID  types.String `tfsdk:"resource_id"`
	AccountID   types.String `tfsdk:"account_id"`
	Region      types.String `tfsdk:"region"`
	DefaultPort types.Int64  `tfsdk:"default_port"`
	VPCNetwork  types.String `tfsdk:"vpc_network"`
}

func NewPortForwardTargetsDataSource() datasource.DataSource {
	return &portForwardTargetsDataSource{}
}

func (d *portForwardTargetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forward_targets"
}

func (d *portForwardTargetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all available port forward targets.",
		Attributes: map[string]schema.Attribute{
			"targets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"arn":          schema.StringAttribute{Computed: true},
						"service":      schema.StringAttribute{Computed: true},
						"resource_id":  schema.StringAttribute{Computed: true},
						"account_id":   schema.StringAttribute{Computed: true},
						"region":       schema.StringAttribute{Computed: true},
						"default_port": schema.Int64Attribute{Computed: true},
						"vpc_network":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *portForwardTargetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*providerdata.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", fmt.Sprintf("Expected *providerdata.ProviderData, got: %T", req.ProviderData))
		return
	}
	d.client = data.Admin
}

func (d *portForwardTargetsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	targets, err := d.client.ListPortForwardTargets(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing port forward targets", err.Error())
		return
	}

	var state portForwardTargetsDataSourceModel
	for _, t := range targets {
		state.Targets = append(state.Targets, portForwardTargetModel{
			ARN:         types.StringValue(t.ARN),
			Service:     types.StringValue(t.Service),
			ResourceID:  types.StringValue(t.ResourceID),
			AccountID:   types.StringValue(t.AccountID),
			Region:      types.StringValue(t.Region),
			DefaultPort: types.Int64Value(int64(t.DefaultPort)),
			VPCNetwork:  types.StringValue(t.VPCNetwork),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
