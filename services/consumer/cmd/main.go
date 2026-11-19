package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/config"
	"github.com/westcrime77/PulseStream/pkg/kafka"
)

func main() {
	fmt.Println("Consumer started")
	defer func() {
		fmt.Println("Consumer finished")
	}()
	ctx := context.Background()
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
	deadlineCtx, _ := context.WithDeadline(ctx,
		time.Now().Add(5*time.Second))
	for {
		var msg kafkago.Message
		msg, err = kafkaReader.Read(deadlineCtx)
		if err != nil {
			log.Fatalf("failed to read kafka message: %v", err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}
