package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Dialer struct {
	dialer *kafka.Dialer
	logger *log.Logger
}

func NewDialer(logger *log.Logger, accessCert, accessKey, caPem string) (*Dialer, error) {
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

	return &Dialer{
		dialer: dialer,
		logger: logger,
	}, nil
}
