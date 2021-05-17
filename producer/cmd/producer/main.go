package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/popsu/go-website-watcher/producer"
)

const (
	accessCert    = "./kafka_access.cert"
	accessKey     = "./kafka_access.key"
	caPem         = "./ca.pem"
	kafkaTopic    = "go-website-watcher"
	wsConfigFile  = "./website_config.txt"
	checkInterval = 30 * time.Second
)

var (
	serviceURI = os.Getenv("KAFKA_SERVICE_URI")
)

func main() {
	wsConfig, err := producer.WebsiteConfigFromFile(wsConfigFile)
	if err != nil {
		log.Fatalf("error reading config file %s", err)
	}

	logger := log.Default()

	svc, err := producer.New(wsConfig, logger, accessCert, accessKey, caPem, kafkaTopic,
		serviceURI, checkInterval)

	if err != nil {
		log.Fatalf("error initializing producer: %s", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go handleInterrupts(ctx, cancel, &wg, logger)

	wg.Add(1)
	go startProducer(ctx, cancel, &wg, svc, logger)

	wg.Wait()
}

func handleInterrupts(ctx context.Context, cancel context.CancelFunc,
	wg *sync.WaitGroup, logger *log.Logger) {
	defer wg.Done()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-c:
			logger.Println("Signal received", sig)
			cancel()
		case <-ctx.Done():
			logger.Println("Will not listen to signals anymore")
			return
		}
	}
}

func startProducer(ctx context.Context, cancel context.CancelFunc,
	wg *sync.WaitGroup, svc *producer.Service, logger *log.Logger) {
	defer wg.Done()

	go func() {
		<-ctx.Done()
		logger.Println("Gracefully shutting down producer")
		svc.Stop()
	}()

	err := svc.Start(ctx)
	if err != nil {
		logger.Println("Error: ", err)
		cancel()
	}
}
