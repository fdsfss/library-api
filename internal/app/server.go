package app

import (
	"database/sql"
	"library-api/internal/handler"
	"library-api/pkg/config"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
)

type server struct {
	fiber.Handler
	app             *fiber.App
	logger          hclog.Logger
	authorHandler   *handler.AuthorHandler
	bookHandler     *handler.BookHandler
	memberHandler   *handler.MemberHandler
	borrowedHandler *handler.BorrowedHandler
	postgres        *sql.DB
}

func Start() {
	s := new(server)

	err := s.generate()
	if err != nil {
		s.logger.Error("error generating server", "error", err)
		os.Exit(1)
	}

	go func() {
		s.logger.Info("starting server...")

		err = s.app.Listen(":" + config.Get().Port)
		if err != nil {
			s.logger.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	s.gracefulShutdown()
}
