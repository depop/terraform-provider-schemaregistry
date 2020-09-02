# Terraform Provider Schemaregistry

## Usage 

```hcl
provider "schemaregistry" {
  registry_host = "http://localhost:8081"
}

resource "registry_subject_schema" "schema" {
  subject       = "schema"
  schema_type   = "AVRO"
  compatibility = "BACKWARD"
  schema        = file("schema.avsc")

  reference {
    name    = "Reference1"
    subject = "event-value"
    version = "1"
  }
}
```

