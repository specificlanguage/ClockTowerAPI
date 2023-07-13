package db

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var GameDB *gorm.DB

func Init() {

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatalf("Could not read environment variables: %s", envErr.Error())
	}

	db, connErr := gorm.Open(postgres.Open(os.Getenv("PG_URL")), &gorm.Config{})
	if connErr != nil {
		log.Fatalf("Could not connect to database: %s", connErr.Error())
	}

	migrateErr := db.AutoMigrate(&Game{}, &GamePlayer{}, &Action{})
	if migrateErr != nil {
		log.Fatalf("Could not migrate database: %s", migrateErr.Error())
	}

	GameDB = db
}
