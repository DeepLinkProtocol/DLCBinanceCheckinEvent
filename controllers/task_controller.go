package controllers

import (
	"go-signin-service/services"
	"go-signin-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// @Summary Get Task Completion Status
// @Description This endpoint checks if the user has completed the sign-in task based on the wallet address.
// @Tags task completion
// @Accept json
// @Produce json
// @Param walletAddress query string true "User Wallet Address"
// @Success 200 {object} utils.Response{data=utils.TaskCompletionData} "Task completion status code="000000" message="success""
// @Failure 400 {object} utils.Response{} "Invalid Argument code="000006""
// @Failure 500 {object} utils.Response{} "System Error code="000002" message="system busy""
// @Router /v1/task/completion [get]
func GetTaskCompletion(c *gin.Context) {
	walletAddress := c.Query("walletAddress")
	if walletAddress == "" {
		utils.JSONResponse(c, http.StatusBadRequest, "000006", "invalid argument: walletAddress is required", nil)
		return
	}

	// Redis Key
	redisKey := services.GetSignInCountKey(walletAddress)

	signInCount, err := services.GetRedisIntValue(redisKey)
	if err != nil {
		if err == redis.Nil {
			signInCount = 0
		} else {
			logrus.Errorln("GetRedisIntValue error: ", err)
			utils.JSONResponse(c, http.StatusInternalServerError, "000002", "system busy", nil)
			return
		}
	}

	completionStatus := signInCount >= 3
	utils.JSONResponse(c, http.StatusOK, "000000", "success", utils.TaskCompletionData{
		Sign: completionStatus,
	})

}
