package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/westcrime77/PulseStream/pkg/config"
	"go.uber.org/fx"
)

type Logger interface {
	Printf(format string, v ...any)
}

func NewServer(lc fx.Lifecycle, log Logger, h http.Handler, cfg config.Http) *http.Server {
	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Printf("http server listening on %s", srv.Addr)

			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("http server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
