package app

import (
	"os"
	"os/signal"
	"syscall"
)

func (s *server) gracefulShutdown() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit
	s.logger.Info("graceful shutdown started")

	err := s.app.Shutdown()
	if err != nil {
		s.logger.Error("error shutting down", "err", err)
	}

	err = s.postgres.Close()
	if err != nil {
		s.logger.Error("error closing postgres", "err", err)
	}

	s.logger.Info("server stopped")

	return
}
