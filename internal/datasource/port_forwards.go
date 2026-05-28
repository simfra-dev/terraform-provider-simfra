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

var _ datasource.DataSource = &portForwardsDataSource{}

type portForwardsDataSource struct {
	client *client.Client
}

type portForwardsDataSourceModel struct {
	PortForwards []portForwardSummaryModel `tfsdk:"port_forwards"`
}

type portForwardSummaryModel struct {
	ID           types.String `tfsdk:"id"`
	TargetARN    types.String `tfsdk:"target_arn"`
	LocalPort    types.Int64  `tfsdk:"local_port"`
	LocalAddress types.String `tfsdk:"local_address"`
	Service      types.String `tfsdk:"service"`
	Status       types.String `tfsdk:"status"`
}

func NewPortForwardsDataSource() datasource.DataSource {
	return &portForwardsDataSource{}
}

func (d *portForwardsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forwards"
}

func (d *portForwardsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all active Simfra port forwards.",
		Attributes: map[string]schema.Attribute{
			"port_forwards": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"target_arn":    schema.StringAttribute{Computed: true},
						"local_port":    schema.Int64Attribute{Computed: true},
						"local_address": schema.StringAttribute{Computed: true},
						"service":       schema.StringAttribute{Computed: true},
						"status":        schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *portForwardsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portForwardsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	forwards, err := d.client.ListPortForwards(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing port forwards", err.Error())
		return
	}

	var state portForwardsDataSourceModel
	for _, f := range forwards {
		state.PortForwards = append(state.PortForwards, portForwardSummaryModel{
			ID:           types.StringValue(f.ID),
			TargetARN:    types.StringValue(f.TargetARN),
			LocalPort:    types.Int64Value(int64(f.LocalPort)),
			LocalAddress: types.StringValue(f.LocalAddress),
			Service:      types.StringValue(f.Service),
			Status:       types.StringValue(f.Status),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
