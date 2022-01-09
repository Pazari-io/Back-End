package handlers

import (
	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/models"

	"github.com/gofiber/fiber/v2"

	"encoding/base64"
	"encoding/json"
)

func DownloadPurchased(c *fiber.Ctx) error {

	// front end should check with smart contract i
	// if the user has purchased the file
	// this also needs to allow download watermaked copy in case of ebook

	c.Accepts("application/json")
	key := c.Query("fileID")

	if key != "" {

		sDec, _ := base64.StdEncoding.DecodeString(key)
		fileName := string(sDec)

		status, err := models.CheckStatusByName(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		original, err := models.GetOrignalFile(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		// might need to add some license files and archive it
		if status == "done" {
			return c.Download(original)

		}

	}
	return nil
}

func DownloadWaterMarked(c *fiber.Ctx) error {

	c.Accepts("application/json")
	key := c.Query("fileID")

	if key != "" {

		sDec, _ := base64.StdEncoding.DecodeString(key)
		fileName := string(sDec)

		taskType, err := models.CheckTypeByName(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}
		if !(taskType == "audio" || taskType == "image" || taskType == "video") {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		status, err := models.CheckStatusByName(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		waterMarkedFile, err := models.GetWaterMarked(fileName, database.DBInstance)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)

		}

		var res models.Results
		err = json.Unmarshal(waterMarkedFile.Results, &res)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		watermarked := "uploads/watermarked/" + res.WaterMaked

		if status == "done" {
			return c.Download(watermarked)

		}

	}

	return nil

}
