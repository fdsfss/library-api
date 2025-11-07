package handler

import (
	"library-api/internal/model"

	"github.com/gofiber/fiber/v2"
)

type borrowedStore interface {
	Create(book *model.Borrowed) error
	Get(id string) ([]model.Book, error)
	Delete(memberId string, bookId string) error
	DeleteList(memberId string, books []string) error
}

func (b *BorrowedHandler) Create(c *fiber.Ctx) error {
	var borrowed model.Borrowed
	err := c.BodyParser(&borrowed)
	if err != nil {
		b.logger.Error("parsing borrowed data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "borrowed book creation failed",
		})
	}

	err = b.store.Create(&borrowed)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "borrowed book creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "borrowed book created",
	})
}

func (b *BorrowedHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	books, err := b.store.Get(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	if len(books) == 0 {
		b.logger.Info("no books found for this member", "id", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "no books found for this member",
		})
	}

	return c.Status(fiber.StatusOK).JSON(books)
}

func (b *BorrowedHandler) Delete(c *fiber.Ctx) error {
	memberId := c.Params("id")
	bookId := c.Params("book_id")
	err := b.store.Delete(memberId, bookId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "borrowed book deleted",
	})
}

func (b *BorrowedHandler) DeleteList(c *fiber.Ctx) error {
	var books []string
	err := c.BodyParser(&books)
	if err != nil {
		b.logger.Error("parsing borrowed data", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "borrowed book delete failed",
		})
	}

	id := c.Params("id")
	err = b.store.DeleteList(id, books)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "borrowed books deleted",
	})
}
