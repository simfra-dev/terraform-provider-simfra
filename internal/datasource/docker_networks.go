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

var _ datasource.DataSource = &dockerNetworksDataSource{}

type dockerNetworksDataSource struct {
	client *client.Client
}

type dockerNetworksDataSourceModel struct {
	Networks []dockerNetworkModel `tfsdk:"networks"`
}

type dockerNetworkModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Driver types.String `tfsdk:"driver"`
	Scope  types.String `tfsdk:"scope"`
}

func NewDockerNetworksDataSource() datasource.DataSource {
	return &dockerNetworksDataSource{}
}

func (d *dockerNetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docker_networks"
}

func (d *dockerNetworksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Simfra-managed Docker networks.",
		Attributes: map[string]schema.Attribute{
			"networks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Computed: true},
						"name":   schema.StringAttribute{Computed: true},
						"driver": schema.StringAttribute{Computed: true},
						"scope":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *dockerNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dockerNetworksDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	networks, err := d.client.ListDockerNetworks(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing Docker networks", err.Error())
		return
	}

	var state dockerNetworksDataSourceModel
	for _, n := range networks {
		state.Networks = append(state.Networks, dockerNetworkModel{
			ID:     types.StringValue(n.ID),
			Name:   types.StringValue(n.Name),
			Driver: types.StringValue(n.Driver),
			Scope:  types.StringValue(n.Scope),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
