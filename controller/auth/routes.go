package auth

import (
	"fmt"
	"kala/exception"
	"kala/model"
	"kala/repository"
	"kala/service"
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
		pass, err := service.Redis_New().Get(fmt.Sprintf("%s:changePassword", user.Email))
		exception.ResponseStatusError_New(err)
		if pass != binding.Password {
			c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("Password not match"))
			return
		}
	}

	if !user.Verified && user.Role == "admin" {
		c.JSON(http.StatusForbidden, "Your account are not verified")
		return
		// exception.ResponseStatusError_New(errors.New("You are not verified yet"))
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
