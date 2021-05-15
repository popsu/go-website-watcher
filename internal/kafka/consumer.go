package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// HandlerFunc is called by the consumer to process received messages
type HandlerFunc func(ctx context.Context, msg string) error

type Consumer struct {
	dialer    *kafka.Dialer
	brokerURI string
	cancel    context.CancelFunc

	handlerFn HandlerFunc
}
