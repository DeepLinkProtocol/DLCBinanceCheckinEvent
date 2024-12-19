package main

import (
	"fmt"
	"go-signin-service/config"
	_ "go-signin-service/docs"
	"go-signin-service/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var port int = 8090

func main() {
	// docs.SwaggerInfo.BasePath = "/api/v1"
	config.InitRedisClient()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.SetupRoutes(r)

	logrus.Infof("server start at port %d", port)
	r.Run(fmt.Sprintf(":%d", port))
}
