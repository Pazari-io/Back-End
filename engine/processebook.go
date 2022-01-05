package engine

import (
	"os"
	"path/filepath"

	"github.com/Pazari-io/Back-End/models"
	"github.com/Pazari-io/Back-End/utils"
	"gorm.io/gorm"
)

func ProcessEbook(fileName string, db *gorm.DB, errChannel chan error) {

	// get em from config later
	pdfCpuPath := "/Users/mac/go/bin/pdfcpu"
	waterMarkFile := "data/pazari-darkest.png"

	extention := filepath.Ext(fileName)
	waterMarkedFileName := utils.ShaHash() + extention

	// Step I: watermark the PDF
	// we can change to position make it user defined text as well
	args := []string{"stamp", "add", "-mode", "image", waterMarkFile, "op:0.15", fileName, "uploads/watermarked/" + waterMarkedFileName}
	_, err := utils.ExecuteCommand(pdfCpuPath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// Step II: encrypt the pdf
	//pdfcpu encrypt -m aes -k 256 -perm none -upw userdecryptkey -opw ownerdecryptkey sample_pdf.pdf sample_encrypted.pdf

	userEncryptionKey := utils.ShaHash()
	ownerEncryptionKey := utils.ShaHash()

	args = []string{"encrypt", "-m", "aes", "-k", "256", "-perm", "none", "-upw", userEncryptionKey, "-opw", ownerEncryptionKey, "uploads/watermarked/" + waterMarkedFileName, "uploads/encrypted/" + waterMarkedFileName}
	_, err = utils.ExecuteCommand(pdfCpuPath, 360, args...)

	if err != nil {
		errChannel <- err
	}
	// Step III: update the database

	// Delete the unecrypted watermarked file
	os.Remove("uploads/watermarked/" + waterMarkedFileName)

	var results models.Results

	results.WaterMaked = waterMarkedFileName
	results.Extention = extention
	results.OwnerEncryptionKey = ownerEncryptionKey
	results.UserEncryptionKey = userEncryptionKey
	results.Encrypted = "uploads/encrypted/" + waterMarkedFileName

	err = models.UpdateTaskResults(fileName, results, db)
	if err != nil {
		errChannel <- err
	}

	errChannel <- nil
}
