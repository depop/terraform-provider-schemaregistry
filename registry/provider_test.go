package registry

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProviderFactories func(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var testAccProviderFunc func() *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"schemaregistry": testAccProvider,
	}
}

// func TestProvider(t *testing.T) {
// 	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
// }

// func TestProvider_impl(t *testing.T) {
// 	var _ terraform.ResourceProvider = Provider()
// }

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REGISTRY_HOST"); v == "" {
		t.Fatal("REGISTRY_HOST must be set for acceptance tests")
	}
}
