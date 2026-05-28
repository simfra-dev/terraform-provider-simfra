package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/simfra-dev/terraform-provider-simfra/internal/provider"
)

//nolint:unused // used by acceptance tests
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"simfra": providerserver.NewProtocol6WithError(provider.New("test")()),
}

//nolint:unused // used by acceptance tests
func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("SIMFRA_ENDPOINT") == "" {
		t.Skip("SIMFRA_ENDPOINT not set, skipping acceptance test")
	}
}
