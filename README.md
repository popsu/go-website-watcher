# Go website watcher

## What

System that monitors website availability and produces metrics that will be stored in PostgreSQL database. Uses Kafka as message broker.

Contains two services: consumer and producer:

### Producer

Periodically checks the target websites and sends the results to Kafka topic.

### Consumer

Consumes the messages from the Kafka topic and writes them into PostgreSQL database.

### Metrics

  | Metric | Type | Description |
  | ------ | ---- | ----------- |
  | created_at         | TIMESTAMPTZ | Time of check |
  | url                | TEXT        | URL of the checked website |
  | regexp_pattern     | TEXT        | Optional regexp pattern to test if the page contents match |
  | regexp_match       | BOOLEAN     | Optional whether the regexp pattern matches |
  | status_code        | SMALLINT    | HTTP status code of the response |
  | timetofirstbyte_ms | SMALLINT    | TimeToFirstByte response time in milliseconds |

Responses over 30 seconds are considered as no response

## Requirements

- Docker Compose
- Kafka and PostgreSQL

If you don't have Kafka and PostgreSQL available, check [Terraform README](./terraform/README.md) to see how to set them up in [Aiven](https://aiven.io/).

## Usage

Easiest way to deploy these is with docker-compose.

You will need some secrets:

- Kafka Access Key (file: kafka_access.key)
- Kafka Access Certificate (file: kafka_access.cert)
- Kafka CA Certificate (file: ca.pem)
- Kafka Service URI (env value: KAFKA_SERVICE_URI)
- PostgreSQL Service URI (env value: POSTGRES_DBURL)

Check [docker-compose.yml file](./docker-compose.yml) to see which secret files and environment values are needed.

After you have set up these files/environment variables:

```bash
docker-compose up
```

## Development

Requirements:

- Go 1.16+
- [(go-)migrate](https://github.com/golang-migrate/migrate) - Install with `make tools`

## Tests

```bash
make test
```

## Links

- [Terraform Aiven provider](https://github.com/aiven/terraform-provider-aiven) check examples folder
- https://help.aiven.io/en/articles/489572-getting-started-with-aiven-for-apache-kafka
- https://github.com/segmentio/kafka-go

## Todo

- [ ] Fix the data race in producer (run the development binary to see it)
- [ ] Tests
- [ ] packages / program structure
- [ ] Add structured logging
- [ ] Add golangci-lint
- [ ] Add ability to use other message brokers
