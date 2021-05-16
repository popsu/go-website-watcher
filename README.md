# Go website watcher

![ci tests](https://github.com/popsu/go-website-watcher/actions/workflows/tests.yml/badge.svg)

## What

System that monitors website availability and produces metrics that will be stored in a PostgreSQL database. Uses Kafka as message broker. Contains consumer and producer services

### Producer

Periodically checks the target websites and sends the results to a Kafka topic.

### Consumer

Consumes the messages from the Kafka topic and writes them into PostgreSQL database.

### Metrics collected

  | Metric | Type | Description |
  | ------ | ---- | ----------- |
  | created_at         | TIMESTAMPTZ | Time of check |
  | url                | TEXT        | URL of the checked website |
  | regexp_pattern     | TEXT        | (Optional) regexp pattern to test if the page contents match |
  | regexp_match       | BOOLEAN     | (Optional) whether the regexp pattern matches page contents |
  | status_code        | SMALLINT    | HTTP status code of the response |
  | timetofirstbyte_ms | SMALLINT    | Time to first byte (TTFB) response time in milliseconds |

Responses over 30 seconds are considered as no response

## Requirements

- Docker Compose
- Kafka and PostgreSQL
- [(go-)migrate](https://github.com/golang-migrate/migrate)
  - if you have Go installed, install with `make tools`
  - otherwise check [migrate CLI](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md) how to install the binary

If you don't have Kafka and PostgreSQL available, check [Terraform README](./terraform/README.md) to see how to easily set them up in [Aiven](https://aiven.io/).

## Usage

Deploy with docker-compose. You will need some secrets:

- Kafka Access Key (file: `kafka_access.key`)
- Kafka Access Certificate (file: `kafka_access.cert`)
- Kafka CA Certificate (file: `ca.pem`)
- Kafka Service URI (env value: `KAFKA_SERVICE_URI`)
- PostgreSQL Service URI (env value: `POSTGRES_DBURL`)

Check [docker-compose.yml file](./docker-compose.yml) to see which secret files and environment values are needed.

After you have set up these files & environment variables, run db-migrations and start the service with:

```bash
make db-migrate-up
docker-compose up
```

## Development

Requirements:

- [Go](https://golang.org/doc/install) 1.16+

## Tests

```bash
make test
```

## Links

- https://help.aiven.io/en/articles/489572-getting-started-with-aiven-for-apache-kafka
- https://github.com/segmentio/kafka-go

## TODO

- [ ] Fix the data race in producer (run the development binary to see it)
- [ ] Tests
- [ ] packages / program structure
- [ ] Build docker images in CI
- [ ] Add structured logging
- [ ] Add golangci-lint
- [ ] Add ability to use other message brokers
