package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func DockerHash(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}
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
		fmt.Printf("%+s\n", body)
		fmt.Printf("%+v\n", data)

		for _, event := range data.Events {
			if event.Action != "push" {
				continue
			}

			var teamID, problemID string

			record := &models.DockerHash{
				Id:         event.Id,
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
