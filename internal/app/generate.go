package app

import (
	"library-api/internal/handler"
	"library-api/internal/store"
	"library-api/pkg/config"
	"library-api/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hashicorp/go-hclog"
)

func (s *server) generate() error {
	s.Handler = logger.New()

	s.app = fiber.New(
		fiber.Config{
			BodyLimit:             20 * 1024 * 1024,
			DisableStartupMessage: true,
		})

	s.logger = hclog.New(&hclog.LoggerOptions{
		JSONEscapeDisabled: true,
		Level:              hclog.Debug,
		JSONFormat:         true,
	})

	s.useMiddleware()

	postgres, err := db.Connect(config.Get().DbConn)
	if err != nil {
		s.logger.Error("database connection failed", "error", err.Error())
		return err
	}
	s.postgres = postgres

	authorStore := store.NewAuthorStore(s.postgres, s.logger)
	authorHandler := handler.NewAuthorHandler(authorStore, s.logger)
	s.authorHandler = authorHandler

	bookStore := store.NewBookStore(s.postgres, s.logger)
	bookHandler := handler.NewBookHandler(bookStore, s.logger)
	s.bookHandler = bookHandler

	memberStore := store.NewMemberStore(s.postgres, s.logger)
	memberHandler := handler.NewMemberHandler(memberStore, s.logger)
	s.memberHandler = memberHandler

	borrowedStore := store.NewBorrowedStore(s.postgres, s.logger)
	borrowedHandler := handler.NewBorrowedHandler(borrowedStore, s.logger)
	s.borrowedHandler = borrowedHandler

	s.router()

	return nil
}
