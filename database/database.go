package database

import (
	"log"

	"github.com/AlexanderMorozov1919/mobileapp/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {

	db, err := config.InitDB()
	if err != nil {
		return err
	}

	DB = db

	// Автомиграция моделей
	err = migrateModels()
	if err != nil {
		return err
	}

	log.Println("Database tables migrated successfully")
	return nil
}

func migrateModels() error {
	return DB.AutoMigrate(
		&models.Doctor{},
		&models.Patient{},
		&models.ContactInfo{},
		&models.Reception{},
		&models.Allergy{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
