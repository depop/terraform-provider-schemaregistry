package registry

import (
	"fmt"
	"testing"

	registry "github.com/dblooman/schema-registry-client/client"
	"github.com/dblooman/schema-registry-client/client/operations"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRegistrySubjectCreate(t *testing.T) {
	var schemaVerification operations.GetSchemaByVersionOK
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "schemaregistry_subject_schema.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubjectCreate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaRegistrySubjectExists(resourceName, &schemaVerification),
					resource.TestCheckResourceAttr(resourceName, "subject", rName),
					resource.TestCheckResourceAttr(resourceName, "compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
				),
			},
		},
	})
}
func TestAccRegistrySubjectReferences(t *testing.T) {
	var schemaVerification operations.GetSchemaByVersionOK
	rName, rName2 := acctest.RandomWithPrefix("tf-acc-test"), acctest.RandomWithPrefix("tf-acc-test2")
	resourceName1 := "schemaregistry_subject_schema.test"
	resourceName2 := "schemaregistry_subject_schema.test2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubjectCreateReferences(rName, rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSchemaRegistrySubjectExists(resourceName1, &schemaVerification),
					testAccCheckSchemaRegistrySubjectExists(resourceName2, &schemaVerification),
					resource.TestCheckResourceAttr(resourceName1, "subject", rName),
					resource.TestCheckResourceAttr(resourceName2, "subject", rName2),
					resource.TestCheckResourceAttr(resourceName1, "compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceName2, "compatibility", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceName1, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName2, "schema_type", "AVRO"),
				),
			},
		},
	})
}

func testAccCheckSchemaRegistrySubjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*registry.Registry)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "schemaregistry_subject_schema" {
			continue
		}

		registryResp := deleteSchemaVersion(conn, rs.Primary.Attributes["subject"])
		if registryResp != nil {
			return fmt.Errorf("Error finding Schema: %s", registryResp)
		}

		if registryResp != nil {
			return fmt.Errorf("Schema (%s) still exists", rs.Primary)
		}
	}

	return nil
}

func testAccCheckSchemaRegistrySubjectExists(resourceName string, schema *operations.GetSchemaByVersionOK) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*registry.Registry)

		schemaResp, err := getSchemaVersion(conn, rs.Primary.Attributes["subject"], "latest")
		if err != nil {
			return fmt.Errorf("Error finding Schema: %s", err)
		}

		if schemaResp == nil {
			return fmt.Errorf("Schema (%s) not found", rs.Primary.ID)
		}

		schema = schemaResp

		return nil
	}
}

func testAccSubjectCreate(subject string) string {
	return fmt.Sprintf(`
		resource "schemaregistry_subject_schema" "test" {
			subject       = "%s"
			schema_type   = "AVRO"
			compatibility = "BACKWARD"
			schema        = <<-EOF
{
    "type": "record",
    "name": "Test",
    "namespace": "com.terraform.model",
    "fields": [
        {
            "name": "eventType",
            "type": "string"
        }
    ]
}
		EOF
		}
	`, subject)
}
func testAccSubjectCreateReferences(case1, case2 string) string {
	return fmt.Sprintf(`
resource "schemaregistry_subject_schema" "test" {
  subject       = "%s"
  schema_type   = "AVRO"
  compatibility = "BACKWARD"
  schema        = <<-EOF
{
    "type": "record",
    "name": "Test",
    "namespace": "com.terraform.model",
    "fields": [
        {
            "name": "eventType",
            "type": "string"
        }
    ]
}
EOF
}

resource "schemaregistry_subject_schema" "test2" {
  subject       = "%s"
  schema_type   = "AVRO"
  compatibility = "BACKWARD"
  reference {
    name    = "Test"
    subject = "${schemaregistry_subject_schema.test.subject}"
    version = schemaregistry_subject_schema.test.version
  }

  schema     = <<-EOF
{
    "type": "record",
    "name": "Test2",
    "namespace": "com.terraform.model",
    "fields": [
        {
            "name": "eventType",
            "type": "string"
        },
        {
            "name": "baseEvent",
            "type": "Test"
        }
    ]
}
EOF
  depends_on = [schemaregistry_subject_schema.test]
}

	`, case1, case2)
}
