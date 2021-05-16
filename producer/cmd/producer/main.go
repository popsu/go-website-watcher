package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
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
	rePattern  = `put"><noscript>Hello, 世界</`
)

var (
	serviceURI = os.Getenv("KAFKA_SERVICE_URI")
)

func main() {
	const websiteURL = "https://golang.org"

	res, err := pingsite(websiteURL, rePattern)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	message, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	sendMessage(message)
}

func pingsite(url, rePattern string) (*model.Message, error) {
	// TODO Add context with timeout
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// https://github.com/tcnksm/go-httpstat/blob/e866bb2744199f5421f2d795b09dc184aac7adcc/_example/main.go
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	message := &model.Message{
		CreatedAt: time.Now(),
		URL:       url,
		// TimeToFirstByte: result.StartTransfer / time.Millisecond,
		// StatusCode:      res.StatusCode,
	}

	if rePattern != "" {
		re := regexp.MustCompile(rePattern)
		matchFound := re.Match(b)

		message.RegexpPattern = sql.NullString{String: rePattern, Valid: true}
		message.RegexpMatch = sql.NullBool{Bool: matchFound, Valid: true}
	}

	// TODO add only if we get the reply
	if true {

		ttfb := result.StartTransfer / time.Millisecond

		message.TimeToFirstByte = &ttfb
		message.StatusCode = sql.NullInt32{Int32: int32(res.StatusCode), Valid: true}
	}

	return message, nil
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

	u, err := uuid.NewV4()
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
