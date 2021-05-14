package main

import (
	"context"
	"log"
	"os"

	"github.com/popsu/go-website-watcher/pkg"
	"github.com/segmentio/kafka-go"
)

const (
	accessCert = "./kafka_access.cert"
	accessKey  = "./kafka_access.key"
	caPem      = "./ca.pem"
	kafkaTopic = "go-website-watcher"
)

var (
	serviceURI = os.Getenv("KAFKA_SERVICE_URI")
)

func main() {
	run()
}

func run() {
	dialer, err := pkg.NewDialer(accessCert, accessKey, caPem)
	if err != nil {
		log.Fatalf("Error initializing dialer :%s", err)
	}

	// producer
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{serviceURI},
		Topic:   kafkaTopic,
		Dialer:  dialer,
	})

	msg := kafka.Message{
		Key:   []byte("TestKey123"),
		Value: []byte("TestValue123"),
	}

	err = producer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatalf("Error sending message: %s", err)
	}

	err = producer.Close()
	if err != nil {
		log.Fatalf("Error closing producer: %s", err)
	}
}
