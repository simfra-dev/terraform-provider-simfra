package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/simfra-dev/terraform-provider-simfra/internal/awsclient"
	"github.com/simfra-dev/terraform-provider-simfra/internal/client"
	simfradatasource "github.com/simfra-dev/terraform-provider-simfra/internal/datasource"
	simfraresource "github.com/simfra-dev/terraform-provider-simfra/internal/resource"
)

var _ provider.Provider = &simfraProvider{}

type simfraProvider struct {
	version string
}

type simfraProviderModel struct {
	Endpoint      types.String `tfsdk:"endpoint"`
	AdminToken    types.String `tfsdk:"admin_token"`
	SkipTLSVerify types.Bool   `tfsdk:"skip_tls_verify"`
	AccessKey     types.String `tfsdk:"access_key"`
	SecretKey     types.String `tfsdk:"secret_key"`
	Region        types.String `tfsdk:"region"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &simfraProvider{version: version}
	}
}

func (p *simfraProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "simfra"
	resp.Version = p.version
}

func (p *simfraProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Simfra. Manages accounts and infrastructure via the admin API, and creates AWS resources with specific IDs via the AWS API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "Simfra server endpoint URL. Can also be set with the SIMFRA_ENDPOINT environment variable.",
				Optional:    true,
			},
			"admin_token": schema.StringAttribute{
				Description: "Admin API bearer token. Can also be set with the SIMFRA_ADMIN_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"skip_tls_verify": schema.BoolAttribute{
				Description: "Skip TLS certificate verification. Can also be set with the SIMFRA_SKIP_TLS_VERIFY environment variable.",
				Optional:    true,
			},
			"access_key": schema.StringAttribute{
				Description: "AWS access key ID for creating AWS resources. Can also be set with the AWS_ACCESS_KEY_ID environment variable.",
				Optional:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "AWS secret access key for creating AWS resources. Can also be set with the AWS_SECRET_ACCESS_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"region": schema.StringAttribute{
				Description: "AWS region for creating AWS resources. Can also be set with the AWS_REGION environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *simfraProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config simfraProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("SIMFRA_ENDPOINT")
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	if endpoint == "" {
		resp.Diagnostics.AddError("Missing endpoint", "The provider requires an endpoint. Set it in the provider block or via the SIMFRA_ENDPOINT environment variable.")
		return
	}

	adminToken := os.Getenv("SIMFRA_ADMIN_TOKEN")
	if !config.AdminToken.IsNull() {
		adminToken = config.AdminToken.ValueString()
	}

	skipTLS := os.Getenv("SIMFRA_SKIP_TLS_VERIFY") == "true"
	if !config.SkipTLSVerify.IsNull() {
		skipTLS = config.SkipTLSVerify.ValueBool()
	}

	data := &ProviderData{
		Admin: client.New(endpoint, adminToken, skipTLS),
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if !config.AccessKey.IsNull() {
		accessKey = config.AccessKey.ValueString()
	}
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if !config.SecretKey.IsNull() {
		secretKey = config.SecretKey.ValueString()
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	if accessKey != "" && secretKey != "" {
		if region == "" {
			region = "us-east-1"
		}
		data.AWS = awsclient.New(awsclient.Config{
			Endpoint:      endpoint,
			AccessKey:     accessKey,
			SecretKey:     secretKey,
			Region:        region,
			SkipTLSVerify: skipTLS,
		})
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *simfraProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		simfraresource.NewAccountResource,
		simfraresource.NewRoute53ZoneResource,
		simfraresource.NewOrganizationResource,
	}
}

func (p *simfraProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		simfradatasource.NewAccountDataSource,
		simfradatasource.NewAccountsDataSource,
		simfradatasource.NewHealthDataSource,
		simfradatasource.NewServicesDataSource,
		simfradatasource.NewStorageSummaryDataSource,
		simfradatasource.NewDockerSummaryDataSource,
		simfradatasource.NewDockerContainersDataSource,
		simfradatasource.NewDockerNetworksDataSource,
		simfradatasource.NewPortForwardDataSource,
		simfradatasource.NewPortForwardsDataSource,
		simfradatasource.NewPortForwardTargetsDataSource,
		simfradatasource.NewSSOSessionsDataSource,
		simfradatasource.NewCAInfoDataSource,
		simfradatasource.NewDNSPortDataSource,
		simfradatasource.NewSMTPPortDataSource,
	}
}
