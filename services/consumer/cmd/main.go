package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/config"
	"github.com/westcrime77/PulseStream/pkg/kafka"
)

func main() {
	fmt.Println("Consumer started")
	defer func() {
		fmt.Println("Consumer finished")
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := config.LoadDotEnv(); err != nil {
		log.Fatalf("failed to load env file: %v", err)
	}

	_, err := config.LoadPostgres()
	if err != nil {
		log.Fatalf("failed to load postgres config: %v", err)
	}

	kafkaCfg, err := config.LoadKafka()
	if err != nil {
		log.Fatalf("failed to load kafka config: %v", err)
	}

	kafkaReader := kafka.NewReader(kafkaCfg)
	defer func() {
		err = kafkaReader.Close()
		if err != nil {
			log.Printf("failed to close kafka reader: %v", err)
		}
	}()
	var message kafkago.Message
	for {
		message, err = kafkaReader.Fetch(ctx)
		if err != nil {
			log.Fatalf("failed to read kafka message: %v", err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		err = kafkaReader.Commit(ctx, message)
		if err != nil {
			log.Fatalf("failed to commit kafka message: %v", err)
		}
	}
}
