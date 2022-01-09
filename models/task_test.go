package models

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	//sqlite mock
	db, err := gorm.Open(sqlite.Open("TaskMock.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Task{})
	return db
}

func TestCreateTaskRecord(t *testing.T) {

	mockTask := Task{Status: "processing", Type: "image", File: "test.jpg"}
	db := initDB()

	task, err := CreateTaskRecord(mockTask, db)
	if err != nil {
		t.Error("Error creating task record")
	}
	if task != "test.jpg" {
		t.Error("Error creating task record")
	}
	err = os.Remove("TaskMock.db")
	if err != nil {
		t.Error("Error deleting mock file")
	}

}
