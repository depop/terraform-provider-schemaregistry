package registry

import (
	"strings"

	registry "github.com/dblooman/schema-registry-client/client"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider -
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"registry_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGISTRY_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGISTRY_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REGISTRY_PASSWORD", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"schemaregistry_subject_schema": resourceRegistrySubjectSchema(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := strings.Split(d.Get("registry_host").(string), "://")

	transport := httptransport.New(url[1], "", []string{url[0]})
	transport.Consumers["application/vnd.schemaregistry+json"] = runtime.JSONConsumer()
	transport.Consumers["application/vnd.schemaregistry.v1+json"] = runtime.JSONConsumer()
	transport.Producers["application/vnd.schemaregistry+json"] = runtime.JSONProducer()
	transport.Producers["application/vnd.schemaregistry.v1+json"] = runtime.JSONProducer()
	transport.Producers["application/json"] = runtime.JSONProducer()
	if username, ok := d.GetOk("username"); ok {
		transport.DefaultAuthentication = newBasicAuth(username.(string), d.Get("password").(string))
	}

	return registry.New(transport, nil), nil
}

func newBasicAuth(username, password string) runtime.ClientAuthInfoWriter {
	return httptransport.BasicAuth(username, password)
}
