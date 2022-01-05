package models

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Results struct {
	WaterMaked         string   `json:"watermaked"`
	Measurements       string   `json:"measurements"`
	Extention          string   `json:"extention"`
	Duration           string   `json:"duration"`
	BPM                string   `json:"bpm"`
	Encrypted          string   `json:"encrypted"` // for PDF and Archive
	UserEncryptionKey  string   `json:"userencryptionkey"`
	OwnerEncryptionKey string   `json:"ownerencryptionkey"`
	ZipFileList        []string `json:"zipfilelist"`
}

type Task struct {
	gorm.Model
	//ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Status  string // "done" | "processing" | "error"
	Type    string // "audio" | "video" | "image" | "archive" | "ebook"
	File    string // File name
	Results datatypes.JSON
}

func CreateTaskRecord(task Task, db *gorm.DB) (string, error) {

	result := db.Create(&task)
	if result.RowsAffected > 0 {
		return task.File, nil
	}
	return "", result.Error
}

func CheckStatusByName(file string, db *gorm.DB) (string, error) {

	var task Task
	result := db.First(&task, "File = ?", file)
	if result.RowsAffected > 0 {
		return task.Status, nil
	}
	return "", result.Error
}

func GetWaterMarked(file string, db *gorm.DB) (Task, error) {

	var task Task
	result := db.First(&task, "File = ?", file)
	if result.RowsAffected > 0 {
		return task, nil
	}
	return Task{}, result.Error
}

func UpdateTaskResults(file string, res Results, db *gorm.DB) error {

	var task Task
	result := db.First(&task, "File = ?", file)
	if result.RowsAffected > 0 {

		json, _ := json.Marshal(res)
		res := datatypes.JSON(json)
		task.Results = res
		task.Status = "done"

		db.Save(&task)

		return nil
	}
	return result.Error
}
