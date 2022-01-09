package engine

import (
	"archive/zip"
	"os"
	"path/filepath"

	// "strconv"
	// "strings"

	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/models"
	"gorm.io/gorm"
)

func ProcessArchive(fileName string, db *gorm.DB, errChannel chan error) {

	// get these from config
	var SevenZipPath = internal.GetKey("SEVEN_ZIP_PATH")

	// Step I: get files in zip
	files, err := getFilesInZip(fileName)
	if err != nil {
		errChannel <- err
	}

	// Step II: extract filees to a temp directory
	RandomHash := internal.ShaHash()
	outPutDir := "-ouploads/temp/" + RandomHash + "/"
	args := []string{"x", fileName, outPutDir}
	_, err = internal.ExecuteCommand(SevenZipPath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// that . is to not create folders
	outPutDirR := "./uploads/temp/" + RandomHash + "/"
	RandomPass := internal.ShaHash()
	extention := filepath.Ext(fileName)
	encryptedFileName := "./uploads/encrypted/" + internal.ShaHash() + ".7z"

	// // Step III: encrypt temp files and create new file
	args = []string{"a", "-t7z", "-m0=lzma2", "-mx=9", "-mfb=64", "-md=32m", "-ms=on", "-mhe=on", "-p" + RandomPass, encryptedFileName, "-r", outPutDirR + "*"}
	_, err = internal.ExecuteCommand(SevenZipPath, 360, args...)

	//Step IV: delete temp files
	os.RemoveAll(outPutDirR)

	if err != nil {
		errChannel <- err
	}
	//fmt.Println(files)

	// Step XX: update the database
	var results models.Results

	results.Encrypted = encryptedFileName
	results.UserEncryptionKey = RandomPass
	results.ZipFileList = files
	results.Extention = extention

	err = models.UpdateTaskResults(fileName, results, db)
	if err != nil {
		errChannel <- err
	}

	errChannel <- nil
}

func getFile(file *zip.File) (string, error) {
	fileread, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileread.Close()
	fileName := file.Name

	if len(fileName) < 1 {

		return "", err
	}

	return fileName, nil
}

func getFilesInZip(fileName string) ([]string, error) {
	read, err := zip.OpenReader(fileName)
	if err != nil {
		//msg := "Failed to open: %s"
		return nil, err
	}
	defer read.Close()

	var allFiles []string
	for _, file := range read.File {
		aFile, err := getFile(file)
		if err != nil {
			//msg := "Failed to list files: %s"
			return nil, err
		}
		allFiles = append(allFiles, aFile)
	}

	return allFiles, nil
}
