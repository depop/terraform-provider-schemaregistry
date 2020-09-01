package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	registry "github.com/dblooman/schema-registry-client/client"
	"github.com/dblooman/schema-registry-client/client/operations"
	"github.com/dblooman/schema-registry-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegistrySubjectSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceSchemaRegistrySubjectCreate,
		Read:   resourceSchemaRegistrySubjectRead,
		Update: resourceSchemaRegistrySubjectUpdate,
		Delete: resourceSchemaRegistrySubjectDelete,

		Schema: map[string]*schema.Schema{
			"subject": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schema": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
			},
			"schema_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AVRO",
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compatibility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reference": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subject": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSchemaRegistrySubjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*registry.Registry)
	ctx := context.Background()

	subject := d.Get("subject").(string)
	schema := d.Get("schema").(string)
	schemaType := d.Get("schema_type").(string)

	var references []*models.SchemaReference
	if v, ok := d.GetOk("reference"); ok {
		references = expandReferences(v.([]interface{}))
	}

	err := registerSchema(client, subject, schema, schemaType, references)
	if err != nil {
		return err
	}

	schemaResp, err := getSchemaVersion(client, subject, "latest")
	if err != nil {
		return err
	}

	_, err = client.Operations.UpdateSubjectLevelConfig(&operations.UpdateSubjectLevelConfigParams{
		Subject: subject,
		Context: ctx,
		Body: &models.ConfigUpdateRequest{
			Compatibility: d.Get("compatibility").(string),
		},
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %v", schemaResp.Payload)

	d.SetId(fmt.Sprintf("%s-%v", subject, schemaResp.Payload.ID))
	d.Set("schema", schemaResp.Payload.Schema)
	d.Set("version", fmt.Sprint(schemaResp.Payload.Version))

	return resourceSchemaRegistrySubjectRead(d, meta)
}

func resourceSchemaRegistrySubjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*registry.Registry)

	subject, err := getSchemaVersion(client, d.Get("subject").(string), d.Get("version").(string))
	if err != nil {
		return err
	}

	err = d.Set("schema", subject.Payload.Schema)
	if err != nil {
		return err
	}

	return nil
}

func resourceSchemaRegistrySubjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*registry.Registry)
	ctx := context.Background()

	subject := d.Get("subject").(string)
	schema := d.Get("schema").(string)
	schemaType := d.Get("schema_type").(string)
	version := d.Get("version").(string)

	var references []*models.SchemaReference
	if v, ok := d.GetOk("reference"); ok {
		references = expandReferences(v.([]interface{}))
	}

	log.Printf("[DEBUG] %v", references)

	compatible, err := testCompatibility(client, subject, schema, schemaType, version, references)
	if err != nil {
		return err
	}
	if !compatible.Payload.IsCompatible {
		return errors.New("Schema is not compatible")
	}

	err = registerSchema(client, subject, schema, schemaType, references)
	if err != nil {
		return err
	}

	schemaResp, err := getSchemaVersion(client, subject, "latest")
	if err != nil {
		return err
	}

	_, err = client.Operations.UpdateSubjectLevelConfig(&operations.UpdateSubjectLevelConfigParams{
		Subject: subject,
		Context: ctx,
		Body: &models.ConfigUpdateRequest{
			Compatibility: d.Get("compatibility").(string),
		},
	})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s-%v", subject, schemaResp.Payload.ID))
	d.Set("schema", schemaResp.Payload.Schema)
	d.Set("version", fmt.Sprint(schemaResp.Payload.Version))

	return nil
}

func resourceSchemaRegistrySubjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*registry.Registry)
	ctx := context.Background()
	subject := d.Get("subject").(string)

	_, err := client.Operations.DeleteSubject(&operations.DeleteSubjectParams{
		Subject: subject,
		Context: ctx,
	})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func suppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func expandReferences(references []interface{}) []*models.SchemaReference {
	schemaReferences := make([]*models.SchemaReference, 0)

	for _, c := range references {
		param := c.(map[string]interface{})
		version, _ := strconv.Atoi(param["version"].(string))
		schemaReference := &models.SchemaReference{
			Name:    param["name"].(string),
			Subject: param["subject"].(string),
			Version: int32(version),
		}
		schemaReferences = append(schemaReferences, schemaReference)
	}

	return schemaReferences
}

func registerSchema(client *registry.Registry, subject, schema, schemaType string, references []*models.SchemaReference) error {
	ctx := context.Background()

	_, err := client.Operations.Register(&operations.RegisterParams{
		Subject: subject,
		Context: ctx,
		Body: &models.RegisterSchemaRequest{
			References: references,
			Schema:     schema,
			SchemaType: schemaType,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func getSchemaVersion(client *registry.Registry, subject, version string) (*operations.GetSchemaByVersionOK, error) {
	ctx := context.Background()

	schemaResp, err := client.Operations.GetSchemaByVersion(&operations.GetSchemaByVersionParams{
		Context: ctx,
		Subject: subject,
		Version: version,
	})
	if err != nil {
		return nil, err
	}
	return schemaResp, nil
}

func testCompatibility(client *registry.Registry, subject, schema, schemaType, version string, references []*models.SchemaReference) (*operations.TestCompatibilityBySubjectNameOK, error) {
	ctx := context.Background()

	compatible, err := client.Operations.TestCompatibilityBySubjectName(&operations.TestCompatibilityBySubjectNameParams{
		Subject: subject,
		Context: ctx,
		Version: version,
		Body: &models.RegisterSchemaRequest{
			References: references,
			Schema:     schema,
			SchemaType: schemaType,
		},
	})
	if err != nil {
		return nil, err
	}
	return compatible, nil
}
