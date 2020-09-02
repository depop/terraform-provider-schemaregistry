package registry

import (
	"fmt"
	"testing"

	"github.com/depop/logentries"
	lexp "github.com/depop/terraform-provider-logentries/logentries/expect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type LogSetResource struct {
	Name string `tfresource:"name"`
}

func TestAccSubjectCreate(t *testing.T) {
	var logSetResource LogSetResource

	subjectName := fmt.Sprintf("terraform-test-%s", acctest.RandString(8))
	testAccSubjectCreate := fmt.Sprintf(`
		resource "schemaregistry_subject_schema" "test_subject" {
		subject       = "%s"
		schema_type   = "AVRO"
		compatibility = "BACKWARD"
		schema        = file("product-like.avsc")
		}
	`, subjectName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLogentriesLogSetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSubjectCreate,
				Check: lexp.TestCheckResourceExpectation(
					"schemaregistry_subject_schema.test_subject",
					&logSetResource,
					testAccCheckLogentriesLogSetExists,
					map[string]lexp.TestExpectValue{
						"name":     lexp.Equals(subjectName),
						"location": lexp.Equals("terraform.io"),
					},
				),
			},
		},
	})
}

func testAccCheckLogentriesLogSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*logentries.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "logentries_logset" {
			continue
		}

		resp, err := client.LogSet.Read(&logentries.LogSetReadRequest{ID: rs.Primary.ID})

		if err == nil {
			return fmt.Errorf("Log set still exists: %#v", resp)
		}
	}

	return nil
}

func testAccCheckLogentriesLogSetExists(resource string, fact interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LogSet Key is set")
		}

		client := testAccProvider.Meta().(*logentries.Client)

		resp, err := client.LogSet.Read(&logentries.LogSetReadRequest{ID: rs.Primary.ID})

		if err != nil {
			return err
		}
		fmt.Println(resp.ID)

		// res := fact.(*LogSetResource)
		// res.Name = resp.Name

		return nil
	}
}
