package api

import (
	"github.com/baazaouihamza/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Hamza",
		LastName:  "Baazaoui",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Single user hamza")
}
