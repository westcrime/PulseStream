package main

import (
	"log"
	nethttp "net/http"
	"os"

	"github.com/westcrime77/PulseStream/pkg/http"
	"github.com/westcrime77/PulseStream/pkg/kafka"
	"github.com/westcrime77/PulseStream/pkg/logger"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/batcher"
	"github.com/westcrime77/PulseStream/services/ingestion/internal/http/handlers"
	"go.uber.org/fx"

	"github.com/westcrime77/PulseStream/pkg/config"
)

func main() {
	wd, _ := os.Getwd()
	log.Println("working dir:", wd)
	if err := config.LoadDotEnv(); err != nil {
		log.Fatalf("failed to load env file: %v", err)
	}

	kafkaCfg, err := config.LoadKafka()
	if err != nil {
		log.Fatalf("failed to load kafka config: %v", err)
	}

	httpCfg, err := config.LoadHttp()
	if err != nil {
		log.Fatalf("failed to load http config: %v", err)
	}

	application := fx.New(
		fx.Supply(kafkaCfg, httpCfg),

		fx.Provide(
			fx.Annotate(
				func() logger.StdLogger { return logger.StdLogger{Logger: log.Default()} },
				fx.As(new(http.Logger)),
				fx.As(new(batcher.Logger)),
				fx.As(new(handlers.Logger)),
			),

			kafka.NewWriter,
			batcher.New,
			handlers.New,

			fx.Annotate(
				func(h *handlers.Handler) nethttp.Handler { return h.Routes() },
			),

			http.NewServer,
		),

		fx.Invoke(func(*nethttp.Server) {}),
	)

	application.Run()
}
