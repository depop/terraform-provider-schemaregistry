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

resource "schemaregistry_subject_schema" "test2" {
  subject       = "testing2"
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
