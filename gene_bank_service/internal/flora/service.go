package flora

import (
	"project_chimera/gene_bank_service/internal/rabbitmq"

	"github.com/gofiber/fiber/v2"
)

// FloraService defines the interface for flora services
type FloraService interface {
	GetFlora(c *fiber.Ctx) error
	PostFlora(c *fiber.Ctx) error
	DeleteFlora(c *fiber.Ctx) error
}

// FloraService is the concrete implementation of FloraService
type floraService struct {
	rmqHandlers *rabbitmq.Handler
}

func NewFloraService(rmqHandlers *rabbitmq.Handler) FloraService {
	return &floraService{rmqHandlers: rmqHandlers}
}

// GetFlora handler for retrieving flora data
func (s *floraService) GetFlora(c *fiber.Ctx) error {
	return nil
}

// PostFlora handler for adding flora data
func (s *floraService) PostFlora(c *fiber.Ctx) error {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"}
	}
	err := s.rmqHandlers.SendAckRequest(c, "addFlora")
	if err != nil {
		return err
	}
	return nil
}

// DeleteFlora handler for deleting flora data
func (s *floraService) DeleteFlora(c *fiber.Ctx) error {
	return nil
}
