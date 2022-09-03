package controller

import (
	"kala/controller/users"
	"kala/model"
	"kala/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(r *gin.Engine) {

	r.Use(CORSMiddleware())

	AuthRoute := r.Group("/")
	AuthRoute.Use(JWTMiddleware())

	users.UserRegisterRoutes(AuthRoute)
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if strings.Contains(auth, "Basic ") {
			auth = strings.Replace(auth, "Basic ", "", 0)
		}

		jwt := util.JWT_NewSignatureOnly()
		jwt, err := jwt.VerifyToken(auth)

		if err != nil {
			c.JSON(http.StatusUnauthorized, model.HTTPResponse_Message(err.Error()))
			return
		}

		c.Set("auth", jwt)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, date")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
