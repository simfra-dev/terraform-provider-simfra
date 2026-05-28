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

var _ datasource.DataSource = &portForwardDataSource{}

type portForwardDataSource struct {
	client *client.Client
}

type portForwardDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	TargetARN    types.String `tfsdk:"target_arn"`
	TargetIP     types.String `tfsdk:"target_ip"`
	TargetPort   types.Int64  `tfsdk:"target_port"`
	LocalPort    types.Int64  `tfsdk:"local_port"`
	LocalAddress types.String `tfsdk:"local_address"`
	VPCNetwork   types.String `tfsdk:"vpc_network"`
	ContainerID  types.String `tfsdk:"container_id"`
	Service      types.String `tfsdk:"service"`
	ResourceID   types.String `tfsdk:"resource_id"`
	AccountID    types.String `tfsdk:"account_id"`
	Region       types.String `tfsdk:"region"`
	CreatedAt    types.String `tfsdk:"created_at"`
	Status       types.String `tfsdk:"status"`
}

func NewPortForwardDataSource() datasource.DataSource {
	return &portForwardDataSource{}
}

func (d *portForwardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forward"
}

func (d *portForwardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a Simfra port forward by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"target_arn":    schema.StringAttribute{Computed: true},
			"target_ip":     schema.StringAttribute{Computed: true},
			"target_port":   schema.Int64Attribute{Computed: true},
			"local_port":    schema.Int64Attribute{Computed: true},
			"local_address": schema.StringAttribute{Computed: true},
			"vpc_network":   schema.StringAttribute{Computed: true},
			"container_id":  schema.StringAttribute{Computed: true},
			"service":       schema.StringAttribute{Computed: true},
			"resource_id":   schema.StringAttribute{Computed: true},
			"account_id":    schema.StringAttribute{Computed: true},
			"region":        schema.StringAttribute{Computed: true},
			"created_at":    schema.StringAttribute{Computed: true},
			"status":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *portForwardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portForwardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state portForwardDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	session, err := d.client.GetPortForward(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading port forward", err.Error())
		return
	}

	state.ID = types.StringValue(session.ID)
	state.TargetARN = types.StringValue(session.TargetARN)
	state.TargetIP = types.StringValue(session.TargetIP)
	state.TargetPort = types.Int64Value(int64(session.TargetPort))
	state.LocalPort = types.Int64Value(int64(session.LocalPort))
	state.LocalAddress = types.StringValue(session.LocalAddress)
	state.VPCNetwork = types.StringValue(session.VPCNetwork)
	state.ContainerID = types.StringValue(session.ContainerID)
	state.Service = types.StringValue(session.Service)
	state.ResourceID = types.StringValue(session.ResourceID)
	state.AccountID = types.StringValue(session.AccountID)
	state.Region = types.StringValue(session.Region)
	state.CreatedAt = types.StringValue(session.CreatedAt.Format("2006-01-02T15:04:05Z"))
	state.Status = types.StringValue(session.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
