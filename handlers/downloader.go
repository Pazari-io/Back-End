package handlers

import (

	// "path/filepath"
	// "strings"

	// "github.com/Pazari-io/Back-End/engine"

	// "github.com/Pazari-io/Back-End/utils"
	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/models"

	"github.com/gofiber/fiber/v2"

	"encoding/base64"
	"encoding/json"
)

func DownloadWaterMarked(c *fiber.Ctx) error {

	c.Accepts("application/json")
	key := c.Query("fileID")

	if key != "" {

		sDec, _ := base64.StdEncoding.DecodeString(key)
		fileName := string(sDec)

		status, err := models.CheckStatusByName(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		zzz, err := models.GetWaterMarked(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		var res models.Results
		err = json.Unmarshal(zzz.Results, &res)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		watermarked := "./uploads/watermarked/" + res.WaterMaked

		if status == "done" {
			return c.Download(watermarked)

		}

	}

	return c.SendString("hi !")

}
