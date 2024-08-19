package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/lemon-mint/envaddr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.eu.org/envloader"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	envloader.LoadEnvFile(".env")

	ln, err := net.Listen("tcp", envaddr.Get(":46211"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	defer ln.Close()

	closeCh := make(chan struct{}, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		close(closeCh)
	}()

	<-closeCh
	log.Info().Msg("shutting down")
}
