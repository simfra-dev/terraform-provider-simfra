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

var _ datasource.DataSource = &caInfoDataSource{}

type caInfoDataSource struct {
	client *client.Client
}

type caInfoDataSourceModel struct {
	Persistent            types.Bool   `tfsdk:"persistent"`
	Directory             types.String `tfsdk:"directory"`
	RootSubject           types.String `tfsdk:"root_subject"`
	RootSerialNumber      types.String `tfsdk:"root_serial_number"`
	RootNotAfter          types.String `tfsdk:"root_not_after"`
	IntermediateSubject   types.String `tfsdk:"intermediate_subject"`
	IntermediateSerial    types.String `tfsdk:"intermediate_serial_number"`
	IntermediateNotAfter  types.String `tfsdk:"intermediate_not_after"`
}

func NewCAInfoDataSource() datasource.DataSource {
	return &caInfoDataSource{}
}

func (d *caInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ca_info"
}

func (d *caInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves Simfra certificate authority information.",
		Attributes: map[string]schema.Attribute{
			"persistent":                schema.BoolAttribute{Computed: true},
			"directory":                 schema.StringAttribute{Computed: true},
			"root_subject":              schema.StringAttribute{Computed: true},
			"root_serial_number":        schema.StringAttribute{Computed: true},
			"root_not_after":            schema.StringAttribute{Computed: true},
			"intermediate_subject":       schema.StringAttribute{Computed: true},
			"intermediate_serial_number": schema.StringAttribute{Computed: true},
			"intermediate_not_after":     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *caInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *caInfoDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	info, err := d.client.CAInfo(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading CA info", err.Error())
		return
	}

	state := caInfoDataSourceModel{
		Persistent:           types.BoolValue(info.Persistent),
		Directory:            types.StringValue(info.Directory),
		RootSubject:          types.StringValue(info.Root.Subject),
		RootSerialNumber:     types.StringValue(info.Root.SerialNumber),
		RootNotAfter:         types.StringValue(info.Root.NotAfter),
		IntermediateSubject:  types.StringValue(info.Intermediate.Subject),
		IntermediateSerial:   types.StringValue(info.Intermediate.SerialNumber),
		IntermediateNotAfter: types.StringValue(info.Intermediate.NotAfter),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
