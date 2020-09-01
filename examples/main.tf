provider "schemaregistry" {
  registry_host = "http://localhost:8081"
}

resource "registry_subject_schema" "daves" {
  subject       = "test1"
  schema_type   = "AVRO"
  compatibility = "BACKWARD"
  schema        = file("product-like.avsc")

  reference {
    name    = "ActivityEvent"
    subject = "activity-base-event-value"
    version = "1"
  }
}

resource "registry_subject_schema" "daves5" {
  subject       = "test2"
  schema_type   = "AVRO"
  compatibility = "BACKWARD"
  schema        = file("product-like.avsc")
}
