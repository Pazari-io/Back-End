package engine

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pazari-io/Back-End/models"
	"github.com/Pazari-io/Back-End/utils"
	"gorm.io/gorm"
)

// get audio duriation
// ffprobe -i Audio_Sample.mp3  -show_entries format=duration -v quiet -of csv="p=0"

// get audio BPM
// aubio tempo Audio_Sample.mp3

// Watermark audio
//ffmpeg -i main.mp3 -filter_complex "amovie=beep.wav:loop=0,asetpts=N/SR/TB,adelay=10s:all=1[beep]; [0][beep]amix=duration=shortest,volume=2"   out.mp3
//ffmpeg -i main.mp3 '-filter_complex', '[0:a]volume=volume=1[aout0];[1:a]volume=volume=2[aout1];[aout1]aloop=loop=-1:size=2e+09,adelay=2000,atrim=start=0:end=2:duration=6[aconcat];[aout0][aconcat]amix=inputs=2:duration=longest:dropout_transition=4 [aout]',

func ProcessAudio(fileName string, db *gorm.DB) error {

	// get these from config
	var aubioPath = "/usr/local/bin/aubio"
	var waterMarkAudio = "data/pazari_watermark_loop.mp3"
	var ffprobePath = "/usr/local/bin/ffprobe"
	var ffmpegPath = "/usr/local/bin/ffmpeg"

	// Step I: get audio information (duration, BPM)

	args := []string{"-i", fileName, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"}
	getAudioDuration, err := utils.ExecuteCommand(ffprobePath, 360, args...)

	getAudioDuration = strings.TrimSpace(getAudioDuration)
	b, _ := strconv.ParseFloat(getAudioDuration, 32)
	duration := int(b)

	log.Println(duration)

	if err != nil {
		return err
	}

	args = []string{"tempo", fileName}
	getAudioBPM, err := utils.ExecuteCommand(aubioPath, 360, args...)

	d := strings.Split(getAudioBPM, " ")
	g, _ := strconv.ParseFloat(d[0], 32)
	bpm := int(g)

	log.Println(bpm)

	if err != nil {
		return err
	}

	// for this version we use a 7 seconds pazari loop might need more adjustments
	// ffmpeg -ar 48000 -t 3 -f s16le -acodec pcm_s16le -i /dev/zero -f mp3 -y silence.mp3
	// ffmpeg -i pazari_audio_watermark.m4a -i silence.mp3 -filter_complex '[0:0] [1:0] concat=n=2:v=0:a=1 [a]' -map '[a]' -ar 48000 -y loop.mp3

	// Step II: watermark audio

	durantionStr := strconv.Itoa(duration)
	extention := filepath.Ext(fileName)
	waterMarkedFileName := utils.ShaHash() + extention
	args = []string{"-i", fileName, "-stream_loop", "-1", "-i", waterMarkAudio, "-filter_complex", "[1:a][0:a]amix", "-t", durantionStr, "-ar", "48000", "-f", "mp3", "-y", "./uploads/watermarked/" + waterMarkedFileName}
	waterMarkAudio, err = utils.ExecuteCommand(ffmpegPath, 600, args...)

	if err != nil {
		return err
	}

	log.Println(waterMarkAudio)

	// Step III: update the database
	var results models.Results

	results.WaterMaked = waterMarkedFileName
	results.Duration = durantionStr
	results.Extention = extention
	results.BPM = strconv.Itoa(bpm)

	err = models.UpdateTaskResults(fileName, results, db)
	if err != nil {
		return err
	}

	return nil

}
