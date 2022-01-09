package engine

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/models"
	"gorm.io/gorm"
)

func ProcessVideo(fileName string, db *gorm.DB, errChannel chan error) {

	// Step I: get Video height and width

	var ffprobePath = internal.GetKey("FFPROBE_PATH")
	var ffmpegPath = internal.GetKey("FFMPEG_PATH")
	var magickPath = internal.GetKey("MAGICK_PATH")

	var waterMarkImage = internal.GetKey("WATER_MARK_IMAGE_VIDEO")

	//var waterMarkImage = "./data/pazari_15_watermark.png" // 15% opacity dark watermark used for video

	//ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 sample_video.mov
	args := []string{"-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", fileName}
	getVideoMeasurement, err := internal.ExecuteCommand(ffprobePath, 360, args...)

	getVideoMeasurement = strings.TrimSpace(getVideoMeasurement)
	// Get half of the image width and height
	measures := strings.Split(getVideoMeasurement, "x")
	weight := measures[0]
	height := measures[1]

	//log.Println(measures)
	halfHeight, _ := strconv.Atoi(height)
	halfHeight = halfHeight / 2
	halfWeight, _ := strconv.Atoi(weight)
	halfWeight = halfWeight / 2

	measureString := strconv.Itoa(halfWeight) + "x" + strconv.Itoa(halfHeight)

	//log.Println(measureString)

	if err != nil {
		errChannel <- err
	}

	extention := filepath.Ext(fileName)
	waterResizedFileName := internal.ShaHash() + ".png"
	waterMarkedFileName := internal.ShaHash() + extention

	// Step II: resize the watermarked image with half of the original size
	args = []string{"convert", waterMarkImage, "-resize", measureString, "uploads/watermarks/" + waterResizedFileName}
	_, err = internal.ExecuteCommand(magickPath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// Step III: do the watermark
	// ffmpeg -i sample_video.mov -i pazari-resized_15.png -filter_complex "overlay=x=(main_w-overlay_w)/2:y=(main_h-overlay_h)/2" output.mp4
	args = []string{"-i", fileName, "-i", "uploads/watermarks/" + waterResizedFileName, "-filter_complex", "overlay=10:main_h-overlay_h-10", "uploads/watermarked/" + waterMarkedFileName}
	_, err = internal.ExecuteCommand(ffmpegPath, 600, args...)

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
