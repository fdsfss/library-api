package handler

import (
	"library-api/internal/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type bookStore interface {
	Create(book *model.Book) error
	Get() ([]model.Book, error)
	Delete(id string) error
	Update(id string, book *model.Book) error
	Exists(id string) error
}

func (b *BookHandler) Create(c *fiber.Ctx) error {
	var book model.Book
	err := c.BodyParser(&book)
	if err != nil {
		b.logger.Error("book body parsing failed for create", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book creation failed",
		})
	}

	book.ID = uuid.New().String()
	err = b.store.Create(&book)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      book.ID,
		"message": "book created",
	})
}

func (b *BookHandler) Get(c *fiber.Ctx) error {
	books, err := b.store.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	if len(books) == 0 {
		b.logger.Info("", "info", "no books found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "no books found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}

func (b *BookHandler) Update(c *fiber.Ctx) error {
	var book model.Book
	err := c.BodyParser(&book)
	if err != nil {
		b.logger.Error("book body parsing failed for update", "id", book.ID, "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book update failed",
		})
	}

	id := c.Params("id")
	err = b.store.Exists(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "book not found",
		})
	}

	err = b.store.Update(id, &book)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "book update failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book updated",
	})
}

func (b *BookHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := b.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "book has related recordings and cannot be deleted",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book deleted",
	})
}
