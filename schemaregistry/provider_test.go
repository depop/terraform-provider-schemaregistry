package registry

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"schemaregistry": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REGISTRY_HOST"); v == "" {
		t.Fatal("REGISTRY_HOST must be set for acceptance tests")
	}
}
