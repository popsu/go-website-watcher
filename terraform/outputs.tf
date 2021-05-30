output "postgresql_service_uri" {
  value     = aiven_pg.postgresql.service_uri
  sensitive = true
}

output "kafka_access_key" {
  value     = aiven_kafka.kafka_service.kafka[0].access_key
  sensitive = true
}

output "kafka_access_cert" {
  value     = aiven_kafka.kafka_service.kafka[0].access_cert
  sensitive = true
}

output "project_ca_cert" {
  value = data.aiven_project.my_project.ca_cert
  sensitive = true
}

output "kafka_service_uri" {
  value     = aiven_kafka.kafka_service.service_uri
  sensitive = true
}
