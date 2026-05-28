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

var _ datasource.DataSource = &ssoSessionsDataSource{}

type ssoSessionsDataSource struct {
	client *client.Client
}

type ssoSessionsDataSourceModel struct {
	Sessions []ssoSessionModel `tfsdk:"sessions"`
}

type ssoSessionModel struct {
	Token     types.String `tfsdk:"token"`
	UserID    types.String `tfsdk:"user_id"`
	UserName  types.String `tfsdk:"user_name"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	CreatedAt types.String `tfsdk:"created_at"`
	Expired   types.Bool   `tfsdk:"expired"`
}

func NewSSOSessionsDataSource() datasource.DataSource {
	return &ssoSessionsDataSource{}
}

func (d *ssoSessionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_sessions"
}

func (d *ssoSessionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all active SSO sessions.",
		Attributes: map[string]schema.Attribute{
			"sessions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"token":      schema.StringAttribute{Computed: true, Sensitive: true},
						"user_id":    schema.StringAttribute{Computed: true},
						"user_name":  schema.StringAttribute{Computed: true},
						"expires_at": schema.StringAttribute{Computed: true},
						"created_at": schema.StringAttribute{Computed: true},
						"expired":    schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ssoSessionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ssoSessionsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	sessions, err := d.client.ListSSOSessions(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing SSO sessions", err.Error())
		return
	}

	var state ssoSessionsDataSourceModel
	for _, s := range sessions {
		state.Sessions = append(state.Sessions, ssoSessionModel{
			Token:     types.StringValue(s.Token),
			UserID:    types.StringValue(s.UserID),
			UserName:  types.StringValue(s.UserName),
			ExpiresAt: types.StringValue(s.ExpiresAt.Format("2006-01-02T15:04:05Z")),
			CreatedAt: types.StringValue(s.CreatedAt.Format("2006-01-02T15:04:05Z")),
			Expired:   types.BoolValue(s.Expired),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
