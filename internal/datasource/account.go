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

var _ datasource.DataSource = &accountDataSource{}

type accountDataSource struct {
	client *client.Client
}

type accountDataSourceModel struct {
	AccountID          types.String `tfsdk:"account_id"`
	RootAccessKeyID    types.String `tfsdk:"root_access_key_id"`
	RootSecretAccessKey types.String `tfsdk:"root_secret_access_key"`
	CreatedAt          types.String `tfsdk:"created_at"`
}

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a Simfra account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The 12-digit AWS account ID.",
				Required:    true,
			},
			"root_access_key_id": schema.StringAttribute{
				Description: "Root user access key ID.",
				Computed:    true,
				Sensitive:   true,
			},
			"root_secret_access_key": schema.StringAttribute{
				Description: "Root user secret access key.",
				Computed:    true,
				Sensitive:   true,
			},
			"created_at": schema.StringAttribute{
				Description: "Account creation timestamp.",
				Computed:    true,
			},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	acct, err := d.client.GetAccount(ctx, state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading account", err.Error())
		return
	}

	state.AccountID = types.StringValue(acct.AccountID)
	state.RootAccessKeyID = types.StringValue(acct.RootAccessKeyID)
	state.RootSecretAccessKey = types.StringValue(acct.RootSecretAccessKey)
	state.CreatedAt = types.StringValue(acct.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
