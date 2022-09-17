package auth

import (
	"kala/model"
	"kala/repository"
	"kala/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userJSONBinding struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterRoutes(c *gin.Engine) {
	c.POST("/auth", login)
	c.GET("/auth", logout)
}

func login(c *gin.Context) {
	var binding userJSONBinding
	timeOffset := c.GetInt("timezone")

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	user, err := repository.User_New().FindUserByEmail(binding.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("Email not found"))
		return
	}

	if !util.Bcrypt_CheckPasswordHash(binding.Password, user.Password) {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("Password not match"))
		return
	}

	jwt := util.JWT_New(user.Name, user.Email, user.Role, timeOffset)
	token, err := jwt.GenerateToken()

	c.JSON(http.StatusOK, model.HTTPResponse_Data(map[string]string{
		"token": token,
		"role":  user.Role,
	}))
}

func logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
