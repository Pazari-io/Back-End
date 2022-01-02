package handlers

import (
	"path/filepath"
	"strings"

	"github.com/Pazari-io/Back-End/database"
	"github.com/Pazari-io/Back-End/engine"
	"github.com/Pazari-io/Back-End/models"
	"github.com/Pazari-io/Back-End/utils"

	"encoding/base64"

	"github.com/gofiber/fiber/v2"
)

/*
keep a huge chunk of extention out for first version
*/
var AudioExtentions = []string{"mp3", "wav"}
var ImageExtentions = []string{"jpg", "png", "jpeg", "psd", "gif", "bmp", "tiff", "webp", "heic", "svg"}
var VideoExtentions = []string{"mp4", "mpg", "mpeg", "avi", "mkv", "webm", "m4v", "mov", "wmv"}
var EbookExtentions = []string{"epub", "mobi", "pdf"}

// game assts / graphics have to be archive now
var ArchiveExtentions = []string{"zip", "rar", "7z"}

// might need a multiple file uploader later
func Uploader(c *fiber.Ctx) error {

	c.Accepts("multipart/form-data")
	file, err := c.FormFile("file")

	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
	}

	// extract and validate extention
	extention := strings.Replace(filepath.Ext(file.Filename), ".", "", -1)

	//log.Println(noDotExtention)

	fileType := "unknown"

	if inSlice(extention, ArchiveExtentions) {
		fileType = "archive"
	} else if inSlice(extention, AudioExtentions) {
		fileType = "audio"
	} else if inSlice(extention, ImageExtentions) {
		fileType = "image"
	} else if inSlice(extention, VideoExtentions) {
		fileType = "video"
	} else if inSlice(extention, EbookExtentions) {
		fileType = "ebook"
	} else {
		fileType = "unknown"
	}

	if fileType == "unknown" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// secure file name to avoid name collision url leakage
	fileName := utils.ShaHash()
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
			go engine.ProcessImage(filePath, database.DBInstance)

			b64ID := base64.StdEncoding.EncodeToString([]byte(fileID))

			return c.JSON(fiber.Map{"taskID": b64ID})

		}

	default:
		c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendString("K")
	//return c.SendStatus(fiber.StatusBadRequest)

}

func inSlice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
