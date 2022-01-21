package handlers

import (
	"path/filepath"
	"strings"

	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/engine"
	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/models"

	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

/*
keep a huge chunk of extention out for first version
*/
var AudioExtentions = []string{"mp3", "wav"}
var ImageExtentions = []string{"jpg", "png", "jpeg", "psd", "gif", "bmp", "tiff", "webp", "heic", "svg"}
var VideoExtentions = []string{"mp4", "mpg", "mpeg", "avi", "mkv", "webm", "m4v", "mov", "wmv"}
var EbookExtentions = []string{"pdf"} // "epub", "mobi",
// game assts / graphics have to be archive now
var ArchiveExtentions = []string{"zip"} // "rar", "7z"

// might need a multiple file uploader later
func Uploader(c *fiber.Ctx) error {

	c.Accepts("multipart/form-data")
	file, err := c.FormFile("file")

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// extract and validate extention
	extention := strings.Replace(filepath.Ext(file.Filename), ".", "", -1)

	fileType := "unknown"

	if internal.InSlice(extention, ArchiveExtentions) {
		fileType = "archive"
	} else if internal.InSlice(extention, AudioExtentions) {
		fileType = "audio"
	} else if internal.InSlice(extention, ImageExtentions) {
		fileType = "image"
	} else if internal.InSlice(extention, VideoExtentions) {
		fileType = "video"
	} else if internal.InSlice(extention, EbookExtentions) {
		fileType = "ebook"
	} else {
		fileType = "unknown"
	}

	if fileType == "unknown" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// secure file name to avoid name collision url leakage
	fileName := internal.ShaHash()
	filePath := "./uploads/original/" + fileName + "." + extention

	c.SaveFile(file, filePath)

	task := models.Task{}
	task.File = filePath
	task.Status = "processing"

	switch fileType {
	case "image":
		{
			task.Type = "image"
			fileID, err := models.CreateTaskRecord(task, database.DBInstance)
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
			}
			errChannel := make(chan error, 1)

			go engine.ProcessImage(filePath, database.DBInstance, errChannel)

			if (<-errChannel) != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})

		}
	case "audio":
		{
			task.Type = "audio"
			fileID, err := models.CreateTaskRecord(task, database.DBInstance)
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
			}
			errChannel := make(chan error, 1)
			go engine.ProcessAudio(filePath, database.DBInstance, errChannel)
			if (<-errChannel) != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})
		}

	case "video":
		{
			task.Type = "video"
			fileID, err := models.CreateTaskRecord(task, database.DBInstance)
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
			}
			errChannel := make(chan error, 1)
			go engine.ProcessVideo(filePath, database.DBInstance, errChannel)
			if (<-errChannel) != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})
		}

	case "ebook":
		{
			task.Type = "ebook"
			fileID, err := models.CreateTaskRecord(task, database.DBInstance)
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
			}
			errChannel := make(chan error, 1)

			go engine.ProcessEbook(filePath, database.DBInstance, errChannel)

			if (<-errChannel) != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})
		}

	case "archive":
		{
			task.Type = "archive"
			fileID, err := models.CreateTaskRecord(task, database.DBInstance)
			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
			}
			errChannel := make(chan error, 1)
			go engine.ProcessArchive(filePath, database.DBInstance, errChannel)
			if (<-errChannel) != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})
		}

	default:
		c.SendStatus(fiber.StatusBadRequest)
	}

	return nil
}
