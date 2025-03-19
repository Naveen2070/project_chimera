package handlers

import (
	"project_chimera/gene_bank_service/internal/services"

	"github.com/gofiber/fiber/v2"
)

func HelloHandler(c *fiber.Ctx) error {
	result := services.GetHelloMessage()

	return c.SendString(result)
}
