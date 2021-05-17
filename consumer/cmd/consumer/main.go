package main

import (
	"context"
	"log"
	"os"

	"github.com/popsu/go-website-watcher/consumer"
)

const (
	accessCert      = "./kafka_access.cert"
	accessKey       = "./kafka_access.key"
	caPem           = "./ca.pem"
	kafkaTopic      = "go-website-watcher"
	consumerGroupID = "my-group-1"
)

var (
	serviceURI = os.Getenv("KAFKA_SERVICE_URI")
	pgURL      = os.Getenv("POSTGRES_DBURL")
)

func main() {
	logger := log.Default()

	svc, err := consumer.New(logger, accessCert, accessKey, caPem, kafkaTopic, serviceURI,
		consumerGroupID, pgURL)
	if err != nil {
		log.Fatalf("Error initializing consumer: %s", err)
	}

	svc.Start(context.Background())
}
