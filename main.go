package main

import (
	"HealthChecker/controller"

	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	controller.Init()
	router := gin.Default()

	v1 := router.Group("/api/v1/urlChecker")
	{
		v1.POST("/", controller.AddURL)
		v1.GET("/", controller.FetchStatus)
		v1.PUT("/:id", controller.UpdateURLData)
	}
	router.Run(":8100")
}
