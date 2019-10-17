package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	models "gitlab.com/CBCTF/bullseye-runner/pkg/master"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func GetRounds(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		rounds := []models.Round{}
		db.Find(&rounds)
		return c.JSON(http.StatusOK, rounds)
	}
}

// GetSchedule returns all schedules currently registered
func GetSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		log.Printf("id:%s", id)
		if id == "" {
			schedules := []models.Schedule{}
			db.Find(&schedules)
			return c.JSON(http.StatusOK, schedules)
		} else {
			schedule := models.Schedule{}
			hit := 0
			db.Where("id = ?", id).Find(&schedule).Count(&hit)
			if hit == 0 {
				return c.JSON(http.StatusNotFound, "schedule not found")
			}
			return c.JSON(http.StatusOK, schedule)
		}
	}
}

// PostSchedule creates new schedule
func PostSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		schedule := models.Schedule{}
		if err := c.Bind(&schedule); err != nil {
			return err
		}
		db.Create(&schedule)
		for t := schedule.StartAt; t.Before(schedule.StopAt); t = t.Add(time.Duration(schedule.Interval) * time.Minute) {
			fmt.Printf("%+v\n", t)
			round := models.Round{
				StartAt:      t,
				Yml:          schedule.Yml,
				FlagTemplate: schedule.FlagTemplate,
				Ntrials:      schedule.Ntrials,
				Timeout:      schedule.Timeout,
				WorkerHosts:  schedule.WorkerHosts,
				CallbackURL:  schedule.CallbackURL,
				ProblemID:    schedule.ProblemID,
				TeamID:       schedule.TeamID,
				Schedule:     schedule,
			}
			db.Create(&round)
		}
		return c.JSON(http.StatusOK, schedule)
	}
}

func DeleteSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		schedule := models.Schedule{}
		hit := 0
		db.Where("id = ?", id).Find(&schedule).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "schedule not found")
		}
		db.Delete(schedule)
		return c.JSON(http.StatusOK, schedule)
	}
}

func GetResults(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		results := []models.Result{}
		db.Find(&results)
		return c.JSON(http.StatusOK, results)
	}
}

func GetWorkerResults(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		workerResults := []models.WorkerResult{}
		db.Find(&workerResults)
		return c.JSON(http.StatusOK, workerResults)
	}
}

func DockerHash(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		records := []models.DockerHash{}
		db.Find(&records)
		return c.JSON(http.StatusOK, records)
	}
}

type Events struct {
	Events []struct {
		Id        string    `json:"id"`
		Timestamp time.Time `json:"timestamp"`
		Action    string    `json:"action"`
		Target    struct {
			MediaType  string `json:"mediaType"`
			Digest     string `json:"digest"`
			Repository string `json:"repository"`
			Tag        string `json:"tag"`
		} `json:"target"`
		Request struct {
			Addr      string `json:"addr"`
			Useragent string `json:"useragent"`
		} `json:"request"`
	} `json:"events"`
}

func Notification(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		data := new(Events)
		if err := json.Unmarshal(body, data); err != nil {
			return err
		}
		for _, event := range data.Events {
			if event.Action != "push" {
				continue
			}

			if event.Target.MediaType != "application/vnd.docker.distribution.manifest.v2+json" {
				continue
			}

			var teamID, problemID string
			if strings.Contains(event.Target.Repository, "/") {
				sep := strings.Split(event.Target.Repository, "/")
				teamID = sep[0]
				problemID = sep[1]
			} else {
				problemID = event.Target.Repository
			}

			record := models.DockerHash{
				Uuid:       event.Id,
				Timestamp:  event.Timestamp,
				Digest:     event.Target.Digest,
				TeamID:     teamID,
				ProblemID:  problemID,
				RemoteAddr: event.Request.Addr,
				UserAgent:  event.Request.Useragent,
			}
			db.Create(&record)
		}
		return c.String(http.StatusOK, "ok")
	}
}
