package main

import (
	"context"
	"fmt"
	"log"

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
	writeMessages(ctx, kafkaWriter, 1000)
}

func writeMessages(ctx context.Context, writer *kafka.Writer, number int64) {
	for i := range number {
		_ = writer.Write(ctx, kafka.Message{
			Value: fmt.Sprintf("value%d", i),
			Key:   fmt.Sprintf("key%d", i),
		})
	}
}
