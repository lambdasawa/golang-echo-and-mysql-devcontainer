package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AccessLog struct {
	ID        uint
	Method    string
	URL       string
	CreatedAt time.Time
}

func main() {
	// https://gorm.io/ja_JP/docs/connecting_to_the_database.html
	dsn := "root:root@tcp(db:3306)/foo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Migrator().CreateTable(&AccessLog{}); err != nil {
		// noop
	}

	e := echo.New()

	// https://echo.labstack.com/middleware/logger/
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		accessLog := AccessLog{
			Method: c.Request().Method,
			URL:    c.Request().RequestURI,
		}
		if err := db.Create(&accessLog).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, accessLog)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
