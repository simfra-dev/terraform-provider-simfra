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

var _ datasource.DataSource = &dockerContainersDataSource{}

type dockerContainersDataSource struct {
	client *client.Client
}

type dockerContainersDataSourceModel struct {
	Containers []dockerContainerModel `tfsdk:"containers"`
}

type dockerContainerModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Image   types.String `tfsdk:"image"`
	State   types.String `tfsdk:"state"`
	Status  types.String `tfsdk:"status"`
	Created types.String `tfsdk:"created"`
}

func NewDockerContainersDataSource() datasource.DataSource {
	return &dockerContainersDataSource{}
}

func (d *dockerContainersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docker_containers"
}

func (d *dockerContainersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Simfra-managed Docker containers.",
		Attributes: map[string]schema.Attribute{
			"containers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true},
						"name":    schema.StringAttribute{Computed: true},
						"image":   schema.StringAttribute{Computed: true},
						"state":   schema.StringAttribute{Computed: true},
						"status":  schema.StringAttribute{Computed: true},
						"created": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *dockerContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dockerContainersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	containers, err := d.client.ListDockerContainers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing Docker containers", err.Error())
		return
	}

	var state dockerContainersDataSourceModel
	for _, c := range containers {
		state.Containers = append(state.Containers, dockerContainerModel{
			ID:      types.StringValue(c.ID),
			Name:    types.StringValue(c.Name),
			Image:   types.StringValue(c.Image),
			State:   types.StringValue(c.State),
			Status:  types.StringValue(c.Status),
			Created: types.StringValue(c.Created),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
