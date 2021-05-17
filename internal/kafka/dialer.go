package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	dialer   *kafka.Dialer
	producer *kafka.Writer
	logger   *log.Logger
}

func New(accessCert, accessKey, caPem, kafkaServiceURI, kafkaTopic string, logger *log.Logger) (*Service, error) {
	dialer, err := NewDialer(accessCert, accessKey, caPem)
	if err != nil {
		return nil, err
	}

	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaServiceURI},
		Topic:   kafkaTopic,
		Dialer:  dialer,
	})

	return &Service{
		dialer:   dialer,
		producer: producer,
		logger:   logger,
	}, nil
}

func (s *Service) Close() error {
	s.logger.Println("Closing kafka producer")
	return s.producer.Close()
}

func (s *Service) SendMessage(message []byte) error {
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

	err = s.producer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	s.logger.Println("Message sent successfully")

	return nil
}

func NewDialer(accessCert, accessKey, caPem string) (*kafka.Dialer, error) {
	keypair, err := tls.LoadX509KeyPair(accessCert, accessKey)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caPem)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, err
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{keypair},
			RootCAs:      caCertPool,
		},
	}

	return dialer, nil
}
