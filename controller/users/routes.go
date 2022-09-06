package users

import (
	"encoding/json"
	"fmt"
	"kala/model"
	"kala/repository"
	"kala/repository/entity"
	"kala/service"
	"kala/util"
	"net/http"
	"strconv"

	ctr "kala/controller/util"

	"github.com/gin-gonic/gin"
)

type createUserJSONBinding struct {
	Email       string `json:"email" gorm:"primaryKey"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Name        string `json:"name"`
}

type selectUserBindingJSON struct {
	Email string `json:"email"`
}

type confirmUserBindingJSON struct {
	selectUserBindingJSON
	RedirectLink string `json:"link"`
}

type passwordBindingJSON struct {
	Token    string `json:"token"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const REDIS_TOKEN_EXP = 60 * util.TOKENEXPIRED

func UserRegisterRoutes(c *gin.RouterGroup, ctx *gin.Engine) {
	c.POST("/user", createUserSalesORAdminByAdmin)
	c.PUT("/user/:id", updateUserSalesORAdminByAdmin)
	c.DELETE("/user/:id", deleteSalesOrAdminByAdmin)
	c.GET("/user", selectUsers)

	ctx.POST("/user/forgot-password", forgotPasswordByEmail)
	ctx.GET("/user/confirm-token", confirmToken)
	ctx.PATCH("/user/change-password", changePassword)
	ctx.POST("/user/confirm-user", sendConfirmAccountToken)
}

func createUserSalesORAdminByAdmin(c *gin.Context) {
	var binding createUserJSONBinding

	_jwt := ctr.GetJWT(c)
	if msg := _jwt.CheckingThisIsAdmin(); msg != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(msg.Error()))
		return
	}

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var users entity.Users = entity.Users{}
	users.Email = binding.Email
	users.Name = binding.Name
	users.PhoneNumber = binding.PhoneNumber
	users.Role = binding.Role

	err = repository.User_New().CreateUser(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func updateUserSalesORAdminByAdmin(c *gin.Context) {
	var binding createUserJSONBinding

	_jwt := ctr.GetJWT(c)

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if binding.Email == _jwt.UserEmail {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("You not allowed update user"))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var users entity.Users = entity.Users{}
	users.Email = binding.Email
	users.Name = binding.Name
	users.PhoneNumber = binding.PhoneNumber
	users.Role = binding.Role
	users.ID = id

	err = repository.User_New().UpdateUser(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func deleteSalesOrAdminByAdmin(c *gin.Context) {

	jwt := ctr.GetJWT(c)
	if msg := jwt.CheckingThisIsAdmin(); msg != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(msg.Error()))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = repository.User_New().DeleteUsers(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func selectUsers(c *gin.Context) {
	jwt := ctr.GetJWT(c)
	if msg := jwt.CheckingThisIsAdmin(); msg != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(msg.Error()))
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	role := c.DefaultQuery("role", "admin")

	users, err := repository.User_New().FindAll(offset, limit, role)

	var jsonMap []map[string]interface{}
	jsonM, err := json.Marshal(users)
	json.Unmarshal(jsonM, &jsonMap)

	for i := 0; i < len(jsonMap); i++ {

		for _, k := range []string{"password", "token", "token_expired", "login_delay"} {
			if _, ok := jsonMap[i][k]; ok {
				delete(jsonMap[i], k)
			}
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.HTTPResponse_Data(jsonMap))

}

func forgotPasswordByEmail(c *gin.Context) {
	var binding confirmUserBindingJSON

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := repository.User_New().FindUserByEmail(binding.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	token, err := util.TokenGenerator(util.TOKEN_CHANGE_PASSWORD, user.Email)
	link := fmt.Sprintf("%s?token=%s", binding.RedirectLink, token)

	err = service.Redis_New().SetWithExp(fmt.Sprintf("%s:changePassword", user.Email), token, REDIS_TOKEN_EXP)

	if err == nil {
		err = service.SMTP_New().SendConfirmMail("Lupa Password", user.Email, link)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.HTTPResponse_Message(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func confirmToken(c *gin.Context) {

	token := c.DefaultQuery("token", "")
	email := c.DefaultQuery("email", "")
	keyword := c.DefaultQuery("keyword", "changePassword")
	test := c.DefaultQuery("testing", "true")

	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("Token or email invalid"))
		return
	}

	err := util.ValidateWithToken(token, email, keyword, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	if test != "true" {
		service.Redis_New().Del(fmt.Sprintf("%s:%s", email, keyword))
	}

	if keyword == "confirmAccount" {
		user, err := repository.User_New().FindUserByEmail(email)
		if err == nil {
			user.Verified = true
			repository.User_New().UpdateUser(user)
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func changePassword(c *gin.Context) {
	var binding passwordBindingJSON

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	err = util.ValidateWithToken(binding.Token, binding.Email, "changePassword", true)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	service.Redis_New().Del(fmt.Sprintf("%s:changePassword", binding.Email))
	user, err := repository.User_New().FindUserByEmail(binding.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("User not found"))
		return
	}
	pswd, err := util.Bcrypt_HashPassword(binding.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.HTTPResponse_Message(err.Error()))
		return
	}
	user.Password = pswd
	err = repository.User_New().UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}

func sendConfirmAccountToken(c *gin.Context) {
	var binding confirmUserBindingJSON

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := repository.User_New().FindUserByEmail(binding.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	token, err := util.TokenGenerator(util.TOKEN_CONFIRM_ACCOUNT, user.Email)
	link := fmt.Sprintf("%s?token=%s&email=%s", binding.RedirectLink, token, user.Email)

	err = service.Redis_New().SetWithExp(fmt.Sprintf("%s:confirmAccount", user.Email), token, REDIS_TOKEN_EXP)

	if err == nil {
		err = service.SMTP_New().SendConfirmMail("Konfirmasi Akun", user.Email, link)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.HTTPResponse_Message(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
