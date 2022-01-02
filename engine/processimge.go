package engine

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pazari-io/Back-End/models"
	"github.com/Pazari-io/Back-End/utils"
	"gorm.io/gorm"
)

func ProcessImage(fileName string, db *gorm.DB) error {

	// get em from config later
	var magickPath = "/usr/local/bin/magick"
	var waterMarkImage = "./data/pazari-full.png"

	// Step I: get image height and width

	extention := filepath.Ext(fileName)
	args := []string{"identify", "-ping", "-format", "%w:%h", fileName}
	getImageSize, err := utils.ExecuteCommand(magickPath, 360, args...)

	if err != nil {
		return err
	}

	// Get half of the image width and height
	measures := strings.Split(getImageSize, ":")
	weight := measures[0]
	height := measures[1]
	halfHeight, _ := strconv.Atoi(height)
	halfHeight = halfHeight / 2
	halfWeight, _ := strconv.Atoi(weight)
	halfWeight = halfWeight / 2

	// two more files
	waterResizedFileName := utils.ShaHash() + extention
	waterMarkedFileName := utils.ShaHash() + extention

	measureString := strconv.Itoa(halfWeight) + "x" + strconv.Itoa(halfHeight)

	// Step II: resize the watermarked image with half of the original size
	args = []string{"convert", waterMarkImage, "-resize", measureString, "./uploads/watermarks/" + waterResizedFileName}
	_, err = utils.ExecuteCommand(magickPath, 360, args...)

	if err != nil {
		return err
	}

	// Step III: do the watermark
	args = []string{"composite", "-dissolve", "15%", "-gravity", "SouthWest", "./uploads/watermarks/" + waterResizedFileName, fileName, "./uploads/watermarked/" + waterMarkedFileName}
	_, err = utils.ExecuteCommand(magickPath, 360, args...)

	if err != nil {
		return err
	}

	// Step IV: update the database
	var results models.Results

	results.WaterMaked = waterMarkedFileName
	results.Measurements = measureString
	results.Extention = extention

	err = models.UpdateTaskResults(fileName, results, db)
	if err != nil {
		return err
	}

	return nil

}
