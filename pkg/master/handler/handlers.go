package handler

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
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
