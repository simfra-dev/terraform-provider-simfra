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

var _ datasource.DataSource = &servicesDataSource{}

type servicesDataSource struct {
	client *client.Client
}

type servicesDataSourceModel struct {
	Services []serviceModel `tfsdk:"services"`
}

type serviceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Protocols   types.List   `tfsdk:"protocols"`
	Operations  types.List   `tfsdk:"operations"`
}

func NewServicesDataSource() datasource.DataSource {
	return &servicesDataSource{}
}

func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all available Simfra services.",
		Attributes: map[string]schema.Attribute{
			"services": schema.ListNestedAttribute{
				Description: "List of services.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"protocols": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"operations": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *servicesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	services, err := d.client.ListServices(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing services", err.Error())
		return
	}

	var state servicesDataSourceModel
	for _, s := range services {
		protocols, diags := types.ListValueFrom(ctx, types.StringType, s.Protocols)
		resp.Diagnostics.Append(diags...)
		operations, diags := types.ListValueFrom(ctx, types.StringType, s.Operations)
		resp.Diagnostics.Append(diags...)

		state.Services = append(state.Services, serviceModel{
			Name:        types.StringValue(s.Name),
			Description: types.StringValue(s.Description),
			Protocols:   protocols,
			Operations:  operations,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
