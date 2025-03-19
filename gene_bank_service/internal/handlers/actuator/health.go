package actuator

import (
	"project_chimera/gene_bank_service/internal/services/actuator"

	"github.com/gofiber/fiber/v2"
)

func HealthHandler(c *fiber.Ctx) error {
	result := actuator.GetHealthMessage()

	return c.SendString(result)
}
