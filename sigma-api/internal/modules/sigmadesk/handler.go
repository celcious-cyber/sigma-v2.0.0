package sigmadesk

import (
	"github.com/gofiber/fiber/v3"
	"sigma-api/internal/core/models"
)

type DeskHandler struct {
	service DeskService
}

func NewDeskHandler(s DeskService) *DeskHandler {
	return &DeskHandler{s}
}

// GetAnnouncements handles listing announcements
func (h *DeskHandler) GetAnnouncements(c fiber.Ctx) error {
	activeOnly := c.Query("active") == "true"
	announcements, err := h.service.GetAnnouncements(activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(announcements)
}

// GetVisitors handles listing visitor logs
func (h *DeskHandler) GetVisitors(c fiber.Ctx) error {
	visitors, err := h.service.GetVisitors()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(visitors)
}

// CreateAnnouncement handles creating a new message
func (h *DeskHandler) CreateAnnouncement(c fiber.Ctx) error {
	var a models.Announcement
	if err := c.Bind().JSON(&a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.service.CreateAnnouncement(a); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Announcement created successfully"})
}
