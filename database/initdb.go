package database

import (
	"github.com/Pazari-io/Back-End/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Download after purchase
// type EncryptedFile struct {
// 	gorm.Model
// 	Key      string
// 	FileName string
// }

// Original file
// type UploadedFile struct {
// 	gorm.Model
// 	Title         string
// 	OrignilName   string
// 	EncryptedName string
// 	// watermaked lower quality copy of the file
// 	WaterMaked       string
// 	IsWaterMarkReady bool
// }

func InitDB() {
	db, err := gorm.Open(sqlite.Open("Pazari.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Task{})

	DBInstance = db

}
