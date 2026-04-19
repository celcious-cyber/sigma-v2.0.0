package sigmalit

import (
	"github.com/gofiber/fiber/v3"
)

type LibraryHandler struct {
	service LibraryService
}

func NewLibraryHandler(s LibraryService) *LibraryHandler {
	return &LibraryHandler{s}
}

func (h *LibraryHandler) ListBooks(c fiber.Ctx) error {
	books, err := h.service.GetBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(books)
}
