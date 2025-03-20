package flora

import (
	"log"
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/pkg/common"

	"github.com/gofiber/fiber/v2"
)

// FloraHandler defines the interface for flora handlers
type FloraHandler interface {
	GetFlora(c *fiber.Ctx) error
	PostFlora(c *fiber.Ctx) error
	DeleteFlora(c *fiber.Ctx) error
	// Add other methods as needed
}

// floraHandler is the concrete implementation of FloraHandler
type floraHandler struct {
	service FloraService
}

// NewFloraHandler creates a new instance of floraHandler with RabbitMQ handler dependencies
func NewFloraHandler(service FloraService) FloraHandler {
	return &floraHandler{service: service}
}

// GetFlora handler for retrieving flora data
func (h *floraHandler) GetFlora(c *fiber.Ctx) error {
	return nil
}

// PostFlora handler for adding flora data
// @Summary Add a flora data to the database
// @Tags Flora
// @Accept json
// @Produce json
// @Param flora body map[string]interface{} true "Flora data"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /flora [post]
func (h *floraHandler) PostFlora(c *fiber.Ctx) error {
	log.Println("PostFlora handler called")
	err := h.service.PostFlora(c)
	log.Println(err)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(common.SuccessResponse{Status: "Flora submitted successfully"})
}

// DeleteFlora handler for deleting flora data
func (h *floraHandler) DeleteFlora(c *fiber.Ctx) error {
	return nil
}

// FloraRouter sets up the routes for flora endpoints
func FloraRouter(router fiber.Router, rmqHandlers *rabbitmq.Handler) {
	service := NewFloraService(rmqHandlers)
	handler := NewFloraHandler(service)

	router.Get("/", handler.GetFlora)
	router.Post("/", handler.PostFlora)
	router.Delete("/", handler.DeleteFlora)
}
