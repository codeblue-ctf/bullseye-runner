package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	models "gitlab.com/CBCTF/bullseye-runner/pkg/master"
)

type Timestamp time.Time

func (t *Timestamp) UnmarshalParam(src string) error {
	ts, err := time.Parse(time.RFC3339, src)
	*t = Timestamp(ts)
	return err
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

// GetSchedule returns all schedules currently registered
func GetSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all schedules
			schedules := []models.Schedule{}
			db.Find(&schedules)
			return c.JSON(http.StatusOK, schedules)
		}

		// return specific schedule
		schedule := models.Schedule{}
		hit := 0
		db.Preload("Rounds").Where("id = ?", id).Find(&schedule).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "schedule not found")
		}
		return c.JSON(http.StatusOK, schedule)
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
		db.Delete(&schedule)
		return c.JSON(http.StatusOK, schedule)
	}
}

func GetRound(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all rounds
			rounds := []models.Round{}
			db.Find(&rounds)
			return c.JSON(http.StatusOK, rounds)
		}
		// return specific round
		round := models.Round{}
		hit := 0
		db.Preload("Results").Where("id = ?", id).Find(&round).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "round not found")
		}
		return c.JSON(http.StatusOK, round)
	}
}

func GetResult(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all results
			results := []models.Result{}
			db.Find(&results)
			return c.JSON(http.StatusOK, results)
		}
		// return specific result
		result := models.Result{}
		hit := 0
		db.Preload("Jobs").Where("id = ?", id).Find(&result).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "result not found")
		}
		return c.JSON(http.StatusOK, result)
	}
}

// DeleteRound cancel running jobs
func DeleteResult(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		result := models.Result{}
		hit := 0
		db.Where("id = ?", id).Find(&result).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "result not found")
		}
		db.Delete(&result)
		return c.JSON(http.StatusOK, result)
	}
}

func GetJob(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all jobs
			jobs := []models.Job{}
			db.Find(&jobs)
			return c.JSON(http.StatusOK, jobs)
		}
		// return specific job
		job := models.Job{}
		hit := 0
		db.Where("id = ?", id).Find(&job).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "job not found")
		}
		return c.JSON(http.StatusOK, job)
	}
}

func DeleteJob(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		job := models.Job{}
		hit := 0
		db.Where("id = ?", id).Find(&job).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "job not found")
		}
		db.Delete(&job)
		return c.JSON(http.StatusOK, job)
	}
}

func Image(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		records := []models.Image{}
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

			image := models.Image{
				UUID:       event.Id,
				Digest:     event.Target.Digest,
				TeamID:     teamID,
				ProblemID:  problemID,
				RemoteAddr: event.Request.Addr,
				UserAgent:  event.Request.Useragent,
			}
			db.Create(&image)
		}
		return c.String(http.StatusOK, "ok")
	}
}
