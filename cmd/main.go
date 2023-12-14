package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"

	"dough-calculator/internal/app"
	"dough-calculator/internal/domain"
)

func main() {
	listener, application := InitApplication()

	server := application.Server()

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("failed to serve")
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel

	log.Info().Msgf("loan service received a shutdown request %v", sig)

	serverTimeoutContext, cancelFunc := context.WithTimeout(
		context.Background(),
		application.Config().Application.Rest.GraceShutdownTimeoutDuration())
	defer func() {
		cancelFunc()
	}()

	if err := server.Shutdown(serverTimeoutContext); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to shutdown")
	}
}

func InitApplication() (net.Listener, domain.Application) {
	initializer := app.NewApplicationInitializer()

	application, err := initializer.Initialize()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize application")
	}

	restConfig := application.Config().Application.Rest

	listener, err := net.Listen("tcp", restConfig.Server)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	return listener, application
}
