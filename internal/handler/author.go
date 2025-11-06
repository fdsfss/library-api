package handler

import (
	"library-api/internal/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type authorStore interface {
	Create(author *model.Author) error
	Get() ([]model.Author, error)
	Exists(id string) error
	Update(id string, author *model.Author) error
	Delete(id string) error
	GetAuthorsBooks(id string) ([]string, error)
}

func (a *AuthorHandler) Create(c *fiber.Ctx) error {
	var author model.Author
	err := c.BodyParser(&author)
	if err != nil {
		a.logger.Error("author body parsing failed for create", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author creation failed",
		})
	}

	author.ID = uuid.New().String()
	err = a.store.Create(&author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "author created",
	})
}

func (a *AuthorHandler) Get(c *fiber.Ctx) error {
	authors, err := a.store.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	if len(authors) == 0 {
		a.logger.Info("", "info", "authors not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "authors not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(authors)
}

func (a *AuthorHandler) Update(c *fiber.Ctx) error {
	var author model.Author

	err := c.BodyParser(&author)
	if err != nil {
		a.logger.Error("author body parsing failed for update", "id", author.ID, "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author update failed",
		})
	}

	id := c.Params("id")

	err = a.store.Exists(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "author not found",
		})
	}

	err = a.store.Update(id, &author)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "author update failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "author updated",
	})
}

func (a *AuthorHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := a.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "author has related recordings and cannot be deleted",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "author deleted",
	})
}

func (a *AuthorHandler) GetAuthorBooks(c *fiber.Ctx) error {
	id := c.Params("id")
	books, err := a.store.GetAuthorsBooks(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	if len(books) == 0 {
		a.logger.Info("", "info", "no authors found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "book not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}
