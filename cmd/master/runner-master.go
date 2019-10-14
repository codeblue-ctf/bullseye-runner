package main

import (
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/CBCTF/bullseye-runner/pkg/master"
	"gitlab.com/CBCTF/bullseye-runner/pkg/master/handler"
)

func initDB(db *gorm.DB) {
	db.AutoMigrate(
		&master.Schedule{},
		&master.Round{},
		&master.Result{},
	)
}

func main() {
	flag.Parse()

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	initDB(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/schedule", handler.GetSchedule(db))
	e.POST("/schedule", handler.PostSchedule(db))

	// notification endpoint for docker-registry
	e.Any("/notification", handler.Notification(db))

	e.Logger.Fatal(e.Start(":8080"))
}
