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

	go func() {
		err := db.AutoMigrate(&GamePlayer{})

		if err != nil {
			log.Fatalf("Could not migrate database: %s", err.Error())
		}
		err2 := db.AutoMigrate(&Game{})
		if err2 != nil {
			log.Fatalf("Could not migrate database: %s", err2.Error())
		}

	}()

	GameDB = db
}
