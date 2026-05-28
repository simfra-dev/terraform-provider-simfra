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

var _ datasource.DataSource = &accountsDataSource{}

type accountsDataSource struct {
	client *client.Client
}

type accountsDataSourceModel struct {
	Accounts []accountSummaryModel `tfsdk:"accounts"`
}

type accountSummaryModel struct {
	AccountID types.String `tfsdk:"account_id"`
	Alias     types.String `tfsdk:"alias"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewAccountsDataSource() datasource.DataSource {
	return &accountsDataSource{}
}

func (d *accountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *accountsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Simfra accounts.",
		Attributes: map[string]schema.Attribute{
			"accounts": schema.ListNestedAttribute{
				Description: "List of accounts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{Computed: true},
						"alias":      schema.StringAttribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *accountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	accounts, err := d.client.ListAccounts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing accounts", err.Error())
		return
	}

	var state accountsDataSourceModel
	for _, a := range accounts {
		state.Accounts = append(state.Accounts, accountSummaryModel{
			AccountID: types.StringValue(a.AccountID),
			Alias:     types.StringValue(a.Alias),
			CreatedAt: types.StringValue(a.CreatedAt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
