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

var _ datasource.DataSource = &dockerSummaryDataSource{}

type dockerSummaryDataSource struct {
	client *client.Client
}

type dockerSummaryDataSourceModel struct {
	ContainersTotal   types.Int64 `tfsdk:"containers_total"`
	ContainersRunning types.Int64 `tfsdk:"containers_running"`
	ContainersStopped types.Int64 `tfsdk:"containers_stopped"`
	Images            types.Int64 `tfsdk:"images"`
	Networks          types.Int64 `tfsdk:"networks"`
	Volumes           types.Int64 `tfsdk:"volumes"`
}

func NewDockerSummaryDataSource() datasource.DataSource {
	return &dockerSummaryDataSource{}
}

func (d *dockerSummaryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docker_summary"
}

func (d *dockerSummaryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Docker resource summary from Simfra.",
		Attributes: map[string]schema.Attribute{
			"containers_total":   schema.Int64Attribute{Computed: true},
			"containers_running": schema.Int64Attribute{Computed: true},
			"containers_stopped": schema.Int64Attribute{Computed: true},
			"images":             schema.Int64Attribute{Computed: true},
			"networks":           schema.Int64Attribute{Computed: true},
			"volumes":            schema.Int64Attribute{Computed: true},
		},
	}
}

func (d *dockerSummaryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dockerSummaryDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	summary, err := d.client.DockerSummary(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Docker summary", err.Error())
		return
	}

	state := dockerSummaryDataSourceModel{
		ContainersTotal:   types.Int64Value(int64(summary.Containers.Total)),
		ContainersRunning: types.Int64Value(int64(summary.Containers.Running)),
		ContainersStopped: types.Int64Value(int64(summary.Containers.Stopped)),
		Images:            types.Int64Value(int64(summary.Images)),
		Networks:          types.Int64Value(int64(summary.Networks)),
		Volumes:           types.Int64Value(int64(summary.Volumes)),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
