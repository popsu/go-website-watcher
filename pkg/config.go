package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

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
