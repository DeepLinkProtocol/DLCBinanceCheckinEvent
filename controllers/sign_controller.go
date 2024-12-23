package controllers

import (
	"fmt"
	"go-signin-service/services"
	"go-signin-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type SignInRequest struct {
	WalletAddress string `json:"walletAddress" binding:"required" example:"0x5de8477A8A47e7F2c5cE05ad4532861a0AaAc909" description:"the wallet address of user to sign in"`
	Signature     string `json:"signature" binding:"required" example:"0x4ae2890fa2206807b6d25039b9b992cdd866989e9b6cc58d5a672b0d1c7e34f0760fe83744acd2595e89bab14fe818699231c302185047edc67582041b0c30e401" description:"the signature signed by the wallet address of user"`
	Data          string `json:"data" binding:"required" example:"hello" description:"used for make signature"`
}

// @Summary User Signin
// @Description Allows a user to sign in within a specific activity period. Signin count will be incremented and stored in Redis.
// @Tags SignIn
// @Accept json
// @Produce json
// @Param request body SignInRequest true "walletAddress: wallet address, signature: signature, data: used for make signature"
// @Success 200 {object} utils.Response{} "code="000000" message="success""
// @Failure 400 {object} utils.Response{} "Invalid Argument code="000006" message="invalid argument""
// @Failure 500 {object} utils.Response{} "System Error code="000002" message="system busy""
// @Router /v1/task/signin [post]
func SignInHandler(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "000006",
			"message": fmt.Sprintf("invalid argument: %s", err),
			"data":    nil,
		})
		return
	}

	walletAddress := req.WalletAddress
	signature := req.Signature
	data := req.Data

	location, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(location)

	activityStartTime := time.Date(2024, 12, 24, 0, 0, 0, 0, location)
	activityEndTime := time.Date(2025, 1, 6, 23, 59, 59, 0, location)

	if currentTime.Before(activityStartTime) || currentTime.After(activityEndTime) {
		utils.JSONResponse(c, http.StatusBadRequest, "000006", "invalid argument: sign-in is not allowed outside of activity time", nil)
		return
	}

	isValid, _ := utils.VerifySignature(data, signature, walletAddress)
	if !isValid {
		utils.JSONResponse(c, http.StatusUnauthorized, "000003", "invalid signature", nil)
		return
	}

	signedToday, err := services.SignedToday(walletAddress)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "000002", "system busy", nil)
		return
	}
	if signedToday {
		utils.JSONResponse(c, http.StatusBadRequest, "000006", "invalid argument: already signed in today", nil)
		return
	}
	if err := services.SignIn(walletAddress); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "000002", "system busy", nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "000000", "success", nil)
}

// @Summary Get signin count
// @Description Allows a user to sign in within a specific activity period. Signin count will be incremented and stored in Redis.
// @Tags SignCount
// @Accept json
// @Produce json
// @Param walletAddress query string true "User Wallet Address"
// @Success 200 {object} utils.Response{data=utils.SignCountData} "Sign Count(1 times per day) code="000000" message="success""
// @Failure 400 {object} utils.Response{} "Invalid Argument code="000006" message="invalid argument""
// @Failure 500 {object} utils.Response{} "System Error code="000002" message="system busy""
// @Router /v1/task/signin [get]
func SignCountHandler(c *gin.Context) {
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

	utils.JSONResponse(c, http.StatusOK, "000000", "success", utils.SignCountData{
		Count: signInCount,
	})
}
