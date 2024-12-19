package middleware

import (
	"go-signin-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("signature")
		if signature == "" {
			utils.JSONResponse(c, http.StatusUnauthorized, "000003", "invalid signature", nil)
			c.Abort()
			return
		}

		isValid, _ := utils.VerifySignature(c.Request.URL.RawQuery, signature, utils.ValidWalletAddress)
		if !isValid {
			utils.JSONResponse(c, http.StatusUnauthorized, "000003", "invalid signature", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func VerifyRSA() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infoln("c.Request.URL.RawQuery: ", c.Request.URL.RawQuery)

		signature := c.GetHeader("signature")
		if signature == "" {
			logrus.Error("signature is empty")
			utils.JSONResponse(c, http.StatusUnauthorized, "000003", "invalid signature", nil)
			c.Abort()
			return
		}

		publicKey, err := utils.GetPublicKey(utils.RSAPublicKey)
		if err != nil {
			logrus.Error("get public key error:", err)
			utils.JSONResponse(c, http.StatusInternalServerError, "000002", "system busy", nil)
			c.Abort()
			return
		}
		logrus.Infoln("singnature: ", signature)
		logrus.Infoln("publicKey: ", utils.RSAPublicKey)

		isValid, err := utils.RSAVerify(c.Request.URL.RawQuery, publicKey, signature)
		if err != nil {
			logrus.Error("verify signature error:", err)
		}
		if !isValid {
			logrus.Error("signature is invalid")
			utils.JSONResponse(c, http.StatusUnauthorized, "000003", "invalid signature", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
