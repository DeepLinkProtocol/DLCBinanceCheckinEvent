package controllers

import (
	"go-signin-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Server Time
// @Description This endpoint returns the current server time in milliseconds since Unix epoch.
// @Tags server time
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Current server time"
// @Failure 500 {object} utils.Response "System busy"
// @Router /v1/time [get]
func GetServerTime(c *gin.Context) {
	currentTime := utils.GetCurrentTimestamp()

	utils.JSONResponse(c, http.StatusOK, "000000", "success", currentTime)
}
