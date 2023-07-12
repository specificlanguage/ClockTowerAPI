package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
	"os"
)

func main() {

	// Load environment vars
	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatalf("Error loading .env file")
		return
	}

	// Load database
	fmt.Println(os.Getenv("PG_URL"))
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("PG_URL"))))
	DB := bun.NewDB(sqldb, pgdialect.New())
	exec, err := DB.Exec("SELECT 1")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	fmt.Print(exec)

	// Initialize webserver
	r := gin.Default()

	// TODO: convert to .env to make sure that we're able to get from multiple sources
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ClockTower Backend",
			"version": "0.0.1",
		})
	})

	// r.POST("/game/create")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
