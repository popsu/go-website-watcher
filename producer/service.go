package producer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/popsu/go-website-watcher/internal/kafka"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/tcnksm/go-httpstat"
)

type Service struct {
	wsConfig      []WebsiteConfig
	logger        *log.Logger
	kafkaProducer *kafka.Producer
	checkInterval time.Duration
	done          chan bool
	ticker        *time.Ticker
}

func New(wsConfig []WebsiteConfig, logger *log.Logger, accessCert, accessKey,
	caPem, kafkaTopic, kafkaServiceURI string, checkInterval time.Duration) (*Service, error) {
	kafkaProducer, err := kafka.NewProducer(accessCert, accessKey, caPem,
		kafkaServiceURI, kafkaTopic, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		wsConfig:      wsConfig,
		logger:        logger,
		kafkaProducer: kafkaProducer,
		checkInterval: checkInterval,
	}, nil
}

func (s *Service) Stop() {
	s.done <- true
	s.ticker.Stop()
	s.logger.Println("Stop called on producer")
}

func (s *Service) Close() {
	s.kafkaProducer.Close()
}

func (s *Service) Start(ctx context.Context) error {
	s.done = make(chan bool)
	s.ticker = time.NewTicker(s.checkInterval)

	defer s.Close()

	for {
		s.checkAllSites()
		select {
		case <-s.done:
			return nil
		case <-s.ticker.C:
			continue
		}
	}
}

func (s *Service) checkAllSites() {
	for _, wc := range s.wsConfig {
		go func(url, re string) {
			err := s.checkAndSendSite(url, re)
			if err != nil {
				s.logger.Printf("Error when checking site %s, err: %s", url, err)
			}
		}(wc.URL, wc.RePattern)
	}
}

func (s *Service) checkAndSendSite(url, re string) error {
	res, err := s.checkSite(url, re)
	if err != nil {
		return err
	}

	message, err := json.Marshal(res)
	if err != nil {
		return err
	}

	s.kafkaProducer.SendMessage(message)

	return nil
}

func (s *Service) checkSite(siteURL, rePattern string) (*model.Message, error) {
	// Time out the request after 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", siteURL, nil)
	if err != nil {
		return nil, err
	}

	var b []byte

	// https://github.com/tcnksm/go-httpstat/blob/e866bb2744199f5421f2d795b09dc184aac7adcc/_example/main.go
	var result httpstat.Result
	ctx = httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)

	// if there are errors, we will write them later to the message payload
	if err == nil {
		b, err = io.ReadAll(res.Body)
		if err == nil {
			res.Body.Close()
		}
	}

	now := time.Now()

	message := &model.Message{
		CreatedAt: &now,
		URL:       &siteURL,
	}

	// Add regexp data if needed
	if rePattern != "" {
		message.RegexpPattern = &rePattern

		if err == nil {
			re := regexp.MustCompile(rePattern)
			matchFound := re.Match(b)
			message.RegexpMatch = &matchFound
		}
	}

	// Add response data
	if err == nil {
		ttfb := result.StartTransfer / time.Millisecond
		statusCode := int32(res.StatusCode)
		message.TimeToFirstByte = &ttfb
		message.StatusCode = &statusCode
	}

	// Add error data
	if err != nil {
		errorString := err.Error()
		message.Error = &errorString
	}

	return message, nil
}
