package engine

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/models"
	"gorm.io/gorm"
)

func ProcessImage(fileName string, db *gorm.DB, errChannel chan error) {

	// get em from config later
	//var magickPath = internal.GetKey("MAGICK_PATH")
	var convertPath = internal.GetKey("CONVERT_PATH")
	var compositePath = internal.GetKey("COMPOSITE_PATH")
	var identifyPath = internal.GetKey("IDENTIFY_PATH")
	var waterMarkImage = internal.GetKey("WATER_MARK_IMAGE")

	// Step I: get image height and width

	extention := filepath.Ext(fileName)
	args := []string{"-ping", "-format", "%w:%h", fileName}
	getImageSize, err := internal.ExecuteCommand(identifyPath, 360, args...)

	if err != nil {
		errChannel <- err
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
	waterResizedFileName := internal.ShaHash() + extention
	waterMarkedFileName := internal.ShaHash() + extention

	measureString := strconv.Itoa(halfWeight) + "x" + strconv.Itoa(halfHeight)

	// Step II: resize the watermarked image with half of the original size
	args = []string{waterMarkImage, "-resize", measureString, "uploads/watermarks/" + waterResizedFileName}
	_, err = internal.ExecuteCommand(convertPath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// Step III: do the watermark
	args = []string{"-dissolve", "15%", "-gravity", "SouthWest", "uploads/watermarks/" + waterResizedFileName, fileName, "uploads/watermarked/" + waterMarkedFileName}
	_, err = internal.ExecuteCommand(compositePath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// Step IV: update the database
	var results models.Results

	results.WaterMaked = waterMarkedFileName
	results.Measurements = measureString
	results.Extention = extention

	err = models.UpdateTaskResults(fileName, results, db)
	if err != nil {
		errChannel <- err
	}

	errChannel <- nil

}
