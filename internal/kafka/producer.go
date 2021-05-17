package kafka

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Dialer
	producer *kafka.Writer
}

func NewProducer(accessCert, accessKey, caPem, kafkaServiceURI, kafkaTopic string, logger *log.Logger) (*Producer, error) {
	dialer, err := NewDialer(logger, accessCert, accessKey, caPem)
	if err != nil {
		return nil, err
	}

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaServiceURI},
		Topic:   kafkaTopic,
		Dialer:  dialer.dialer,
	})

	return &Producer{
		Dialer:   *dialer,
		producer: producer,
	}, nil
}

func (p *Producer) Close() error {
	p.logger.Println("Closing kafka producer")
	return p.producer.Close()
}

func (p *Producer) SendMessage(message []byte) error {
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}

	key, _ := u.MarshalText()
	// (The error returned should always be nil, so redundant check)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   key,
		Value: message,
	}

	err = p.producer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	p.logger.Println("Message sent successfully")

	return nil
}
