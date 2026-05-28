package resource_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/simfra-dev/terraform-provider-simfra/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"simfra": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("SIMFRA_ENDPOINT") == "" {
		t.Skip("SIMFRA_ENDPOINT not set, skipping acceptance test")
	}
}

func TestAccSimfraAccount_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfig("111222333444"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simfra_account.test", "account_id", "111222333444"),
					resource.TestCheckResourceAttrSet("simfra_account.test", "root_access_key_id"),
					resource.TestCheckResourceAttrSet("simfra_account.test", "root_secret_access_key"),
					resource.TestCheckResourceAttrSet("simfra_account.test", "created_at"),
				),
			},
			{
				ResourceName:            "simfra_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bootstrap", "bootstrap_region", "availability_zones", "vpc_cidr"},
			},
		},
	})
}

func TestAccSimfraAccount_withBootstrap(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfigWithBootstrap("222333444555"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simfra_account.test", "account_id", "222333444555"),
					resource.TestCheckResourceAttr("simfra_account.test", "bootstrap", "standard"),
					resource.TestCheckResourceAttr("simfra_account.test", "bootstrap_region", "us-east-1"),
					resource.TestCheckResourceAttrSet("simfra_account.test", "root_access_key_id"),
				),
			},
		},
	})
}

func testAccAccountConfig(accountID string) string {
	return fmt.Sprintf(`
provider "simfra" {
  endpoint = %[1]q
}

resource "simfra_account" "test" {
  account_id = %[2]q
}
`, os.Getenv("SIMFRA_ENDPOINT"), accountID)
}

func testAccAccountConfigWithBootstrap(accountID string) string {
	return fmt.Sprintf(`
provider "simfra" {
  endpoint = %[1]q
}

resource "simfra_account" "test" {
  account_id       = %[2]q
  bootstrap        = "standard"
  bootstrap_region = "us-east-1"
}
`, os.Getenv("SIMFRA_ENDPOINT"), accountID)
}
