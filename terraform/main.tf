# postgresql
resource "aiven_pg" "postgresql" {
  project      = var.aiven_project_name
  service_name = "go-website-watcher-pq"
  cloud_name   = "google-europe-north1"
  plan         = "hobbyist"

  # Set this to true in production to avoid accidental deletions
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

# We need the project to get CA Cert that Kafka clients require
data "aiven_project" "my_project" {
  project = var.aiven_project_name
}
