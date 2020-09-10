# Terraform Provider Schemaregistry

## Usage 

```hcl
provider "schemaregistry" {
  registry_host = "http://localhost:8081"
}

resource "schemaregistry_subject_schema" "test" {
  subject       = "testing1"
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
```

