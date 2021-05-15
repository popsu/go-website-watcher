package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	gwwkafka "github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/segmentio/kafka-go"
	"github.com/tcnksm/go-httpstat"
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
	const websiteURL = "https://golang.org"

	res, err := pingsite(websiteURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	d := model.Message{
		CreatedAt:    time.Now(),
		URL:          websiteURL,
		TimeToFirstByte: res.StartTransfer / time.Millisecond,
	}

	message, err := json.Marshal(&d)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	sendMessage(message)
}

func pingsite(url string) (*httpstat.Result, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		return nil, err
	}
	res.Body.Close()
	result.End(time.Now())

	return &result, nil
}

func sendMessage(message []byte) error {
	dialer, err := gwwkafka.NewDialer(accessCert, accessKey, caPem)
	if err != nil {
		// log.Fatalf("Error initializing dialer :%s", err)
		return err
	}

	// producer
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{serviceURI},
		Topic:   kafkaTopic,
		Dialer:  dialer,
	})

	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	// the error returned is always nil
	key, _ := u.MarshalText()

	msg := kafka.Message{
		Key:   key,
		Value: message,
	}

	err = producer.WriteMessages(context.Background(), msg)
	if err != nil {
		// log.Fatalf("Error sending message: %s", err)
		return err
	}
	log.Println("message sent successfully")

	err = producer.Close()
	if err != nil {
		// log.Fatalf("Error closing producer: %s", err)
		return err
	}

	return nil
}
