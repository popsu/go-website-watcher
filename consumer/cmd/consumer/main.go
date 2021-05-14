package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/popsu/go-website-watcher/pkg"
	"github.com/segmentio/kafka-go"
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
)

func main() {
	run()
}

func run() {
	dialer, err := pkg.NewDialer(accessCert, accessKey, caPem)
	if err != nil {
		log.Fatalf("Error initializing dialer: %s", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{serviceURI},
		GroupID:  consumerGroupID,
		Topic:    kafkaTopic,
		Dialer:   dialer,
		MinBytes: 10_000,     // 10kB
		MaxBytes: 10_000_000, // 10MB,
	})

	for {
		time.Sleep(time.Second * 1)
		fmt.Println("Polling for messages...")

		// Note this call blocks until a message is available
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error reading message: %s", err)
		}

		fmt.Printf("Message at topic: %v partition: %v offset: %v %s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
