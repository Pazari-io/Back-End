package handlers

import (
	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/models"

	"github.com/gofiber/fiber/v2"

	"encoding/base64"
)

func TaskStatus(c *fiber.Ctx) error {

	c.Accepts("application/json")
	taskID := c.Query("taskID")

	if taskID != "" {

		sDec, _ := base64.StdEncoding.DecodeString(taskID)
		fileName := string(sDec)

		status, err := models.CheckStatusByName(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}
		return c.JSON(fiber.Map{"status": status})

	}

	return nil

}
