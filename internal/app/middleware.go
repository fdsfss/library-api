package app

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *server) useMiddleware() {
	s.app.Use(cors.New(
		cors.Config{
			AllowMethods: "GET,POST,DELETE,PATCH",
			AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
			MaxAge:       120,
		}),
	)

	prometheus := fiberprometheus.New("library-api")
	prometheus.RegisterAt(s.app, "/metrics")
	s.app.Use(prometheus.Middleware)

	s.app.Use(s.selectiveLogging)
}

func (s *server) selectiveLogging(c *fiber.Ctx) error {
	if c.Path() == "/healthz" || c.Path() == "/metrics" {
		return c.Next()
	}

	return s.Handler(c)
}

func (s *server) healthcheck(c *fiber.Ctx) error {
	err := s.postgres.Ping()
	if err != nil {
		s.logger.Error("error pinging database", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error pinging database"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}
