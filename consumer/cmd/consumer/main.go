package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	gwwkafka "github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/popsu/go-website-watcher/internal/persistence"
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

type Consumer struct {
	store  *persistence.PostgresStore
	logger *log.Logger
}

func NewConsumer(dburl string) (*Consumer, error) {
	logger := log.Default()

	store, err := persistence.NewPostgresStore(dburl, logger)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		store:  store,
		logger: logger,
	}, nil
}

func (c *Consumer) writeToDBMsg(msg *model.Message) error {
	return c.store.InsertMessage(msg)
}

func (c *Consumer) writeKafkaMessageToDB(m *kafka.Message) error {
	msg := &model.Message{}

	err := json.Unmarshal(m.Value, &msg)
	if err != nil {
		return err
	}

	uuid, err := uuid.FromString(string(m.Key))
	if err != nil {
		return err
	}

	msg.ID = uuid

	return c.writeToDBMsg(msg)
}

func run() {
	dialer, err := gwwkafka.NewDialer(accessCert, accessKey, caPem)
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

	dburl := os.Getenv("POSTGRES_DBURL")

	c, err := NewConsumer(dburl)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

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

		c.writeKafkaMessageToDB(&m)
	}
}
