package providerdata

import (
	"github.com/simfra-dev/terraform-provider-simfra/internal/awsclient"
	"github.com/simfra-dev/terraform-provider-simfra/internal/client"
)

// ProviderData holds all clients created during provider configuration.
// Passed to resources and data sources via ProviderData/DataSourceData.
type ProviderData struct {
	Admin *client.Client
	AWS   *awsclient.Client
}
