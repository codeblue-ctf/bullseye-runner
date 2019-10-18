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
		&master.Job{},
		&master.Image{},
	)
}

func main() {
	flag.Parse()

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	initDB(db)

	go master.RunScheduler(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)

	e.GET("/schedule", handler.GetSchedule(db))
	e.GET("/schedule/:id", handler.GetSchedule(db))
	e.POST("/schedule", handler.PostSchedule(db))
	e.DELETE("/schedule/:id", handler.DeleteSchedule(db))

	e.GET("/round", handler.GetRound(db))
	e.GET("/round/:id", handler.GetRound(db))

	e.GET("/result", handler.GetResult(db))
	e.GET("/result/:id", handler.GetResult(db))
	e.DELETE("/result/:id", handler.DeleteResult(db))

	e.GET("/job", handler.GetJob(db))
	e.GET("/job/:id", handler.GetJob(db))

	e.GET("/image", handler.Image(db))

	// notification endpoint for docker-registry
	e.Any("/notification", handler.Notification(db))

	e.Logger.Fatal(e.Start(":8080"))
}
