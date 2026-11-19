package main

import (
	"context"
	"fmt"
	"log"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/westcrime77/PulseStream/pkg/config"
	"github.com/westcrime77/PulseStream/pkg/kafka"
)

func main() {
	fmt.Println("Producer started")
	defer func() {
		fmt.Println("Producer finished")
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

	kafkaWriter := kafka.NewWriter(kafkaCfg)
	defer func() {
		err = kafkaWriter.Close()
		if err != nil {
			log.Printf("failed to close kafka writer: %v", err)
		}
	}()
	err = writeMessages(ctx, kafkaWriter, 1000)
	if err != nil {
		log.Fatalf("failed to produce messages: %v", err)
	}
}

func writeMessages(ctx context.Context, writer *kafka.Writer, number int64) error {
	messages := make([]kafkago.Message, number)
	for i := range number {
		messages = append(messages, kafkago.Message{
			Value: []byte(fmt.Sprintf("value%d", i)),
			Key:   []byte(fmt.Sprintf("key%d", i)),
		})
	}
	return writer.Write(ctx, messages)
}
