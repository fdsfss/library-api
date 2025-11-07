package handler

import (
	"library-api/internal/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type memberStore interface {
	Create(member *model.Member) error
	Get() ([]model.Member, error)
	Exists(id string) error
	Update(id string, member *model.Member) error
	Delete(id string) error
}

func (m *MemberHandler) Create(c *fiber.Ctx) error {
	var member model.Member
	err := c.BodyParser(&member)
	if err != nil {
		m.logger.Error("member body parsing failed for create", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member creation failed",
		})
	}

	member.ID = uuid.New().String()
	err = m.store.Create(&member)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "member created",
	})
}

func (m *MemberHandler) Get(c *fiber.Ctx) error {
	members, err := m.store.Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	if len(members) == 0 {
		m.logger.Info("no members found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "no members found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(members)
}

func (m *MemberHandler) Update(c *fiber.Ctx) error {
	var member model.Member

	err := c.BodyParser(&member)
	if err != nil {
		m.logger.Error("member body parsing failed for update", "id", member.ID, "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member update failed",
		})
	}

	id := c.Params("id")

	err = m.store.Exists(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "member not found",
		})
	}

	err = m.store.Update(id, &member)
	if err != nil {
		m.logger.Error("member update failed", "id", id, "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member update failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member updated",
	})
}

func (m *MemberHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := m.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "member still has books, all books must be returned",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member deleted",
	})
}
