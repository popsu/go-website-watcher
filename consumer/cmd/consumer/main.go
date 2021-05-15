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
	"github.com/jackc/pgx/v4"
	gwwkafka "github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
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

		err = writeToDB(&m)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
}

func writeToDB(m *kafka.Message) error {
	ctx := context.Background()

	dburl := os.Getenv("GWW_DBURL")

	uuid, err := uuid.FromString(string(m.Key))
	if err != nil {
		return err
	}

	d := model.Message{}

	err = json.Unmarshal(m.Value, &d)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, dburl)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	ct, err := conn.Exec(ctx, sqlQuery,
		uuid,
		d.CreatedAt,
		d.URL,
		d.RegexpPattern,
		d.RegexpMatch,
		d.StatusCode,
		d.TimeToFirstByte,
	)

	if err != nil {
		return err
	}

	fmt.Println("Rows affected: ", ct.RowsAffected())

	return nil
}
