package handler

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	master "gitlab.com/CBCTF/bullseye-runner/pkg/master"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

// GetSchedule returns all schedules currently registered
func GetSchedule(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all schedules
			schedules := []master.Schedule{}
			db.Find(&schedules)
			return c.JSON(http.StatusOK, schedules)
		}

		// return specific schedule
		schedule := master.Schedule{}
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
		schedule := master.Schedule{}
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
		schedule := master.Schedule{}
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
			rounds := []master.Round{}
			db.Find(&rounds)
			return c.JSON(http.StatusOK, rounds)
		}
		// return specific round
		round := master.Round{}
		hit := 0
		db.Preload("Results").Where("id = ?", id).Find(&round).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "round not found")
		}
		return c.JSON(http.StatusOK, round)
	}
}

// PostRound is for re-evaluation by hand
func PostRound(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		round := master.Round{}
		if err := c.Bind(&round); err != nil {
			return err
		}
		db.Create(&round)
		return c.JSON(http.StatusOK, round)
	}
}

func DeleteRound(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		round := master.Round{}
		hit := 0
		db.Preload("Results").Find(&round, "id = ?", id).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "round not found")
		}
		db.Delete(&round)
		return c.JSON(http.StatusOK, round)
	}
}

func GetResult(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			// return all results
			results := []master.Result{}
			db.Find(&results)
			return c.JSON(http.StatusOK, results)
		}
		// return specific result
		result := master.Result{}
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
		result := master.Result{}
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
			jobs := []master.Job{}
			db.Find(&jobs)
			return c.JSON(http.StatusOK, jobs)
		}
		// return specific job
		job := master.Job{}
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
		job := master.Job{}
		hit := 0
		db.Where("id = ?", id).Find(&job).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "job not found")
		}
		db.Delete(&job)
		return c.JSON(http.StatusOK, job)
	}
}

func GetJobCapture(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.Param("uuid")
		path := fmt.Sprintf("%s/%s.mp4", master.X11CapPrefix, uuid)
		return c.File(path)
	}
}

func GetSampleCaptureByRoundID(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		round := master.Round{}
		hit := 0
		db.Find(&round, "id = ?", id).Count(&hit)
		if hit == 0 {
			return c.JSON(http.StatusNotFound, "round not found")
		}

		results := []master.Result{}
		db.Model(&round).Related(&results)
		if len(results) == 0 {
			return c.JSON(http.StatusNotFound, "result not found")
		}

		result := results[0]
		jobs := []master.Job{}
		db.Model(&result).Related(&jobs)

		succeeded := []string{}
		failed := []string{}
		for _, job := range jobs {
			if job.Succeeded {
				succeeded = append(succeeded, job.UUID)
			} else {
				failed = append(failed, job.UUID)
			}
		}
		return c.JSON(http.StatusOK, struct {
			Succeeded []string `json:"succeeded"`
			Failed    []string `json:"failed"`
		}{
			Succeeded: succeeded,
			Failed:    failed,
		})
	}
}

func ListRunning(c echo.Context) error {
	return c.JSON(http.StatusOK, master.CancelMgr.Keys())
}

func Image(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		records := []master.Image{}
		db.Find(&records)
		return c.JSON(http.StatusOK, records)
	}
}
