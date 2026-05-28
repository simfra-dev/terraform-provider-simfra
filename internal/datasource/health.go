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

var _ datasource.DataSource = &healthDataSource{}

type healthDataSource struct {
	client *client.Client
}

type healthDataSourceModel struct {
	Ready types.Bool `tfsdk:"ready"`
}

func NewHealthDataSource() datasource.DataSource {
	return &healthDataSource{}
}

func (d *healthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_health"
}

func (d *healthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Checks Simfra server health status.",
		Attributes: map[string]schema.Attribute{
			"ready": schema.BoolAttribute{
				Description: "Whether the Simfra server is ready.",
				Computed:    true,
			},
		},
	}
}

func (d *healthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *healthDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	ready, err := d.client.Health(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error checking health", err.Error())
		return
	}

	state := healthDataSourceModel{
		Ready: types.BoolValue(ready),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
