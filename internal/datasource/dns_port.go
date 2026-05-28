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

var _ datasource.DataSource = &dnsPortDataSource{}

type dnsPortDataSource struct {
	client *client.Client
}

type dnsPortDataSourceModel struct {
	AccountID types.String `tfsdk:"account_id"`
	Port      types.Int64  `tfsdk:"port"`
}

func NewDNSPortDataSource() datasource.DataSource {
	return &dnsPortDataSource{}
}

func (d *dnsPortDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_port"
}

func (d *dnsPortDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the DNS server port for a Simfra account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "AWS account ID.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "DNS server port.",
				Computed:    true,
			},
		},
	}
}

func (d *dnsPortDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsPortDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	port, err := d.client.DNSPort(ctx, state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS port", err.Error())
		return
	}

	state.Port = types.Int64Value(int64(port))
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
