package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	// sql query embedded from file
	_ "embed"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	gwwkafka "github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/popsu/go-website-watcher/internal/psql"
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

//go:embed sql/insert.sql
var sqlQuery string

func main() {
	run()
}

type Consumer struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewConsumer(dburl string) (*Consumer, error) {
	db, err := psql.New(dburl)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		db:     db,
		logger: log.Default(),
	}, nil
}

func (c *Consumer) writeToDBMsg(msg *model.Message) error {
	log.Default()
	ct, err := c.db.Exec(sqlQuery,
		msg.ID,
		msg.CreatedAt,
		msg.URL,
		msg.RegexpPattern,
		msg.RegexpMatch,
		msg.StatusCode,
		msg.TimeToFirstByte,
	)

	if err != nil {
		return err
	}

	n, err := ct.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("Rows affected: ", n)

	return nil
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

	dburl := os.Getenv("GWW_DBURL")

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
