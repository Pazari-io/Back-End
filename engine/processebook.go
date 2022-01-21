package engine

import (
	"path/filepath"

	"github.com/Pazari-io/Back-End/internal"
	"github.com/Pazari-io/Back-End/models"
	"gorm.io/gorm"
)

func ProcessEbook(fileName string, db *gorm.DB, errChannel chan error) {

	// get em from config later
	pdfCpuPath := internal.GetKey("PDF_CPU_PATH")
	waterMarkFile := internal.GetKey("WATER_MARK_IMAGE_DARK")

	extention := filepath.Ext(fileName)
	waterMarkedFileName := internal.ShaHash() + extention

	// Step I: watermark the PDF
	// we can change to position make it user defined text as well
	args := []string{"stamp", "add", "-mode", "image", waterMarkFile, "op:0.15", fileName, "uploads/watermarked/" + waterMarkedFileName}
	_, err := internal.ExecuteCommand(pdfCpuPath, 360, args...)

	if err != nil {
		errChannel <- err
	}

	// Step II: encrypt the pdf
	//pdfcpu encrypt -m aes -k 256 -perm none -upw userdecryptkey -opw ownerdecryptkey sample_pdf.pdf sample_encrypted.pdf

	userEncryptionKey := internal.ShaHash()
	ownerEncryptionKey := internal.ShaHash()

	args = []string{"encrypt", "-m", "aes", "-k", "256", "-perm", "none", "-upw", userEncryptionKey, "-opw", ownerEncryptionKey, "uploads/watermarked/" + waterMarkedFileName, "uploads/encrypted/" + waterMarkedFileName}
	_, err = internal.ExecuteCommand(pdfCpuPath, 360, args...)

	if err != nil {
		errChannel <- err
	}
	// Step III: update the database

	// Delete the unecrypted watermarked file
	// os.Remove("uploads/watermarked/" + waterMarkedFileName)

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
