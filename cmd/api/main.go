// Command api is the flight-meta search service entrypoint.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flightmeta/internal/combiner"
	"flightmeta/internal/config"
	"flightmeta/internal/httpapi"
	"flightmeta/internal/search"
	"flightmeta/internal/sources"
	"flightmeta/internal/sources/mockleg"
	"flightmeta/internal/visa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	var srcs []sources.Adapter
	if cfg.EnableMock {
		// Own combiner over a mock per-leg price source. A real price feed
		// (e.g. Travelpayouts) implements sources.LegSource and replaces mockleg
		// here without touching the combiner or the rest of the pipeline.
		srcs = append(srcs, combiner.New(mockleg.New(), nil))
	}
	if len(srcs) == 0 {
		log.Error("no data sources configured; set FM_ENABLE_MOCK=true or add a real leg source")
		os.Exit(1)
	}

	resolver, err := visa.Load()
	if err != nil {
		log.Error("failed to load transit-visa data", "err", err)
		os.Exit(1)
	}

	orch := search.New(log, cfg.SourceTimeout, resolver, srcs...)
	handler := httpapi.New(orch, log, cfg.CORSOrigin)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("listening", "addr", cfg.Addr, "mock", cfg.EnableMock)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("stopped")
}
