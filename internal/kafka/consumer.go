package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofrs/uuid"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Dialer
	reader *kafka.Reader
}

func NewConsumer(accessCert, accessKey, caPem, kafkaTopic, kafkaServiceURI, consumerGroupID string,
	logger *log.Logger) (*Consumer, error) {
	dialer, err := NewDialer(logger, accessCert, accessKey, caPem)
	if err != nil {
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaServiceURI},
		GroupID: consumerGroupID,
		Topic:   kafkaTopic,
		Dialer:  dialer.dialer,
		// Values from https://github.com/segmentio/kafka-go/tree/f0d3749a707d1a63b089082f8a1aa9980566427d#consumer-groups
		MinBytes: 10_000,     // 10kB
		MaxBytes: 10_000_000, // 10MB,
	})

	return &Consumer{
		Dialer: *dialer,
		reader: reader,
	}, nil
}

func (c *Consumer) Close() error {
	c.logger.Println("Closing kafka reader")
	return c.reader.Close()
}

func (c *Consumer) ReadMessage(ctx context.Context) (*model.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)

	if err != nil {
		return nil, err
	}

	return c.kafkaMessageToModelMessage(&msg)
}

func (c *Consumer) kafkaMessageToModelMessage(m *kafka.Message) (*model.Message, error) {
	msg := &model.Message{}

	err := json.Unmarshal(m.Value, &msg)
	if err != nil {
		return nil, err
	}

	uuid, err := uuid.FromString(string(m.Key))
	if err != nil {
		return nil, err
	}

	msg.ID = &uuid

	return msg, nil
}
