package router

import (
	"go-signin-service/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * 3600,
	}))

	// Task Completion API with signature validation
	r.POST("/v1/task/signin", controllers.SignInHandler)

	r.GET("/v1/Time", controllers.GetServerTime)
	r.GET("/v1/task/completion", controllers.GetTaskCompletion)
	r.GET("/v1/task/signin", controllers.SignCountHandler)

}
