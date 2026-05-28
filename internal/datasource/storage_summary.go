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

var _ datasource.DataSource = &storageSummaryDataSource{}

type storageSummaryDataSource struct {
	client *client.Client
}

type storageSummaryDataSourceModel struct {
	Enabled  types.Bool                  `tfsdk:"enabled"`
	DataDir  types.String                `tfsdk:"data_dir"`
	DBPath   types.String                `tfsdk:"db_path"`
	DBSize   types.Int64                 `tfsdk:"db_size"`
	Total    types.Int64                 `tfsdk:"total"`
	Services []storageServiceInfoModel   `tfsdk:"services"`
}

type storageServiceInfoModel struct {
	Service types.String `tfsdk:"service"`
	Size    types.Int64  `tfsdk:"size"`
}

func NewStorageSummaryDataSource() datasource.DataSource {
	return &storageSummaryDataSource{}
}

func (d *storageSummaryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_summary"
}

func (d *storageSummaryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Simfra storage and persistence summary.",
		Attributes: map[string]schema.Attribute{
			"enabled":  schema.BoolAttribute{Computed: true},
			"data_dir": schema.StringAttribute{Computed: true},
			"db_path":  schema.StringAttribute{Computed: true},
			"db_size":  schema.Int64Attribute{Computed: true},
			"total":    schema.Int64Attribute{Computed: true},
			"services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service": schema.StringAttribute{Computed: true},
						"size":    schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *storageSummaryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *storageSummaryDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	summary, err := d.client.StorageSummary(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading storage summary", err.Error())
		return
	}

	state := storageSummaryDataSourceModel{
		Enabled: types.BoolValue(summary.Enabled),
		DataDir: types.StringValue(summary.DataDir),
		DBPath:  types.StringValue(summary.DBPath),
		DBSize:  types.Int64Value(summary.DBSize),
		Total:   types.Int64Value(summary.Total),
	}
	for _, s := range summary.Services {
		state.Services = append(state.Services, storageServiceInfoModel{
			Service: types.StringValue(s.Service),
			Size:    types.Int64Value(s.Size),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
