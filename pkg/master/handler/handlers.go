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

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func GetSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
func PostSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
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

			record := &models.DockerHash{
				Uuid:       event.Id,
				Timestamp:  event.Timestamp,
				Digest:     event.Target.Digest,
				TeamID:     teamID,
				ProblemID:  problemID,
				RemoteAddr: event.Request.Addr,
				UserAgent:  event.Request.Useragent,
			}
			db.Create(record)
		}
		return c.String(http.StatusOK, "ok")
	}
}
