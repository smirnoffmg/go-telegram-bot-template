package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"github.com/smirnoffmg/partner/app"
)

const cancelTimeout = 3 * time.Second

func main() {
	sentryDsn := os.Getenv("SENTRY_DSN")

	if sentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
		}); err != nil {
			log.Error().Err(err).Msg("Cannot init sentry")
		} else {
			defer sentry.Flush(cancelTimeout)
		}
	}

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Error().Err(err).Msg("Run")

		defer os.Exit(1)
	}
}

func run(ctx context.Context) error {
	app, err := app.NewApp()
	if err != nil {
		log.Error().Err(err).Msg("Cannot create app")

		return err
	}

	if err := app.Start(); err != nil {
		log.Error().Err(err).Msg("Problem with bot")

		return err
	}

	<-ctx.Done()

	log.Info().Msg("Shutting down (waiting 3 seconds)...")

	<-time.After(cancelTimeout)

	return nil
}
