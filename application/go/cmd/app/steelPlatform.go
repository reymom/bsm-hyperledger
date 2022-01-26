package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/reymom/bsm-hyperledger/application/go/cmd/app/config"
	"github.com/reymom/bsm-hyperledger/application/go/internal/routes"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("------- Starting Steel Platform App, Version: %s, Build Date: %s -------")

	//app routing
	conf, e := config.GenerateConfig()
	if e != nil {
		log.Err(e).Msg("Error while generating configuration")
		os.Exit(1)
	}

	handler, e := routes.GenerateRoutes(conf)
	if e != nil {
		log.Err(e).Msg("Error while creating routes")
		os.Exit(1)
	}

	appServer := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go startServer(appServer)

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	<-signalChannel
	log.Trace().Msg("Shutting down app")

	e = appServer.Shutdown(context.Background())
	if e != nil {
		log.Err(e).Msg("Error while shutting down the app listener")
	}

	os.Exit(0)
}

func startServer(server *http.Server) {

	if err := server.ListenAndServe(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

}
