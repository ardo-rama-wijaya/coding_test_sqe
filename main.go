package main

import (
	"crypto/subtle"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	dataRequest struct {
		Data string
		Id   string
	}
)

func main() {
	e := echo.New()
	e.POST("/users", saveData)
	e.GET("/users/:id", getData)
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.Logger.Fatal(e.Start(":1323"))
}

func saveData(fc echo.Context) error {
	//panggi.akadol@smma.id
	data := fc.QueryParam("data")
	dsn := "host=localhost user=postgres password=root dbname=sun_mandara port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fc.JSON(http.StatusBadGateway, err.Error())
	}
	err = db.Table("data").Create(dataRequest{Data: data, Id: strconv.Itoa(rand.Intn(100))}).Error
	if err != nil {
		fc.JSON(http.StatusBadGateway, err.Error())
	}
	return fc.JSON(http.StatusOK, err)
}

func getData(fc echo.Context) error {
	dsn := "host=localhost user=postgres password=root dbname=sun_mandara port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fc.JSON(http.StatusBadGateway, err.Error())
	}
	var res dataRequest
	err = db.Table("data").Where("id = ?", fc.Param("id")).Find(&res).Error
	if err != nil {
		fc.JSON(http.StatusBadGateway, err.Error())
	}
	return fc.JSON(http.StatusOK, res)
}
