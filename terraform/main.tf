# postgresql
resource "aiven_pg" "postgresql" {
  project      = var.aiven_project_name
  service_name = "go-website-watcher-pq"
  cloud_name   = "google-europe-north1"
  plan         = "hobbyist"

  termination_protection = false
}

# kafka
resource "aiven_kafka" "kafka_service" {
  project      = var.aiven_project_name
  cloud_name   = "google-europe-north1"
  plan         = "startup-2"
  service_name = "go-website-watcher-kafkasvc"
}

resource "aiven_kafka_topic" "go_website_watcher" {
  project      = var.aiven_project_name
  service_name = aiven_kafka.kafka_service.service_name
  topic_name   = "go-website-watcher"
  partitions   = 5
  replication  = 3
}

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

output "kafka_service_uri" {
  value     = aiven_kafka.kafka_service.service_uri
  sensitive = true
}
