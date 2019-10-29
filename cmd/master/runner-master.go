package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/CBCTF/bullseye-runner/pkg/master"
	"gitlab.com/CBCTF/bullseye-runner/pkg/master/handler"
)

var (
	DbDialect = getenv("DB_DIALECT", "sqlite3")
	DbConnect = getenv("DB_CONNECT", "test.db")
	Port      = getenv("PORT", ":8080")

	debug = flag.Bool("debug", false, "enable debug")
)

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func initDB(db *gorm.DB) {
	db.AutoMigrate(
		&master.Schedule{},
		&master.Round{},
		&master.Result{},
		&master.Job{},
	)
}

func main() {
	flag.Parse()

	db, err := gorm.Open(DbDialect, DbConnect)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	master.InitLogger(*debug)
	initDB(db)

	master.InitScheduler()
	go master.RunScheduler(db)
	go master.RunUpdater(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)

	e.GET("/schedule", handler.GetSchedule(db))
	e.GET("/schedule/:id", handler.GetSchedule(db))
	e.POST("/schedule", handler.PostSchedule(db))
	e.DELETE("/schedule/:id", handler.DeleteSchedule(db))

	e.GET("/round", handler.GetRound(db))
	e.POST("/round", handler.PostRound(db))
	e.GET("/round/:id", handler.GetRound(db))
	e.GET("/round/capture/:id", handler.GetSampleCaptureByRoundID(db))
	e.DELETE("/round/:id", handler.DeleteRound(db))

	e.GET("/result", handler.GetResult(db))
	e.GET("/result/:id", handler.GetResult(db))
	e.DELETE("/result/:id", handler.DeleteResult(db))

	e.GET("/job", handler.GetJob(db))
	e.GET("/job/:id", handler.GetJob(db))
	e.DELETE("/job/:id", handler.DeleteJob(db))
	e.GET("/job/capture/:uuid", handler.GetJobCapture(db))

	e.GET("/running", handler.ListRunning)

	e.GET("/image", handler.Image(db))

	// pprof
	if *debug {
		log.Printf("pprof enabled")
		pprofGroup := e.Group("/debug/pprof")
		pprofGroup.Any("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
		pprofGroup.Any("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		pprofGroup.Any("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
		pprofGroup.Any("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
		pprofGroup.Any("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	}

	e.Logger.Fatal(e.Start(Port))
}
