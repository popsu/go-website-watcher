package consumer

import (
	"context"
	"log"

	"github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/popsu/go-website-watcher/internal/persistence"
)

type Service struct {
	kafkaConsumer *kafka.Consumer
	store         *persistence.PostgresStore
	logger        *log.Logger
}

func New(logger *log.Logger, accessCert, accessKey, caPem, kafkaTopic, kafkaServiceURI string,
	consumerGroupID, postgresURI string) (*Service, error) {

	store, err := persistence.NewPostgresStore(postgresURI, logger)
	if err != nil {
		return nil, err
	}

	kafkaConsumer, err := kafka.NewConsumer(accessCert, accessKey, caPem, kafkaTopic, kafkaServiceURI, consumerGroupID, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		kafkaConsumer: kafkaConsumer,
		store:         store,
		logger:        logger,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	for {
		s.logger.Println("Polling for messages...")
		// Note this call blocks until a message is available
		m, err := s.kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			s.logger.Printf("Error reading message: %s", err)
		}

		s.writeToDBMsg(m)
	}
}

func (s *Service) writeToDBMsg(msg *model.Message) {
	err := s.store.InsertMessage(msg)
	if err != nil {
		s.logger.Println("Error writing to db:", err)
	}
}
