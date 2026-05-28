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

var _ datasource.DataSource = &smtpPortDataSource{}

type smtpPortDataSource struct {
	client *client.Client
}

type smtpPortDataSourceModel struct {
	AccountID types.String `tfsdk:"account_id"`
	Port      types.Int64  `tfsdk:"port"`
}

func NewSMTPPortDataSource() datasource.DataSource {
	return &smtpPortDataSource{}
}

func (d *smtpPortDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_port"
}

func (d *smtpPortDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the SMTP relay port for a Simfra account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "AWS account ID.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "SMTP relay port.",
				Computed:    true,
			},
		},
	}
}

func (d *smtpPortDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *smtpPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state smtpPortDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, err := d.client.SMTPPort(ctx, state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SMTP port", err.Error())
		return
	}

	state.Port = types.Int64Value(int64(port))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
