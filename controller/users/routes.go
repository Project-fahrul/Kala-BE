package users

import (
	"encoding/json"
	"fmt"
	"kala/exception"
	"kala/model"
	"kala/repository"
	"kala/repository/entity"
	"kala/service"
	"kala/util"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	Email string `json:"email" binding:"required"`
}

type confirmUserBindingJSON struct {
	selectUserBindingJSON
	RedirectLink string `json:"link"`
}

type passwordBindingJSON struct {
	selectUserBindingJSON
	Password string `json:"password" binding:"required"`
}

type registerSalesJsonBinding struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone"`
}

type changePasswordAdminBinding struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

const REDIS_TOKEN_EXP = 60 * util.TOKENEXPIRED

func UserRegisterRoutes(c *gin.RouterGroup, ctx *gin.Engine) {
	c.POST("/user", createUserSalesORAdminByAdmin)
	c.PUT("/user/:id", updateUserSalesORAdminByAdmin)
	c.DELETE("/user/:id", deleteSalesOrAdminByAdmin)
	c.GET("/user", selectUsers)
	c.GET("/me", me)
	c.GET("/user/:id", getUser)

	ctx.POST("/user/forgot-password", forgotPasswordByEmail)
	// ctx.GET("/user/confirm-token", confirmToken)
	ctx.PATCH("/user/change-password", changePassword)
	// ctx.POST("/user/confirm-user", sendConfirmAccountToken)
	ctx.POST("/user/email", userByEmail)
	ctx.POST("/user/registration", registerSales)
	c.GET("/user/sales", allSales)
	c.GET("/user/sales-not-verified", listSalesNotVerified)
	c.GET("user/verified/:id/:action", userVerified)

	c.PATCH("/user/changePassword/self", changePasswordUser)
}

func changePasswordUser(c *gin.Context) {
	var pswd changePasswordAdminBinding
	err := c.ShouldBindJSON(&pswd)

	exception.ResponseStatusError_New(err)
	_jwt := ctr.GetJWT(c)
	user, err := repository.User_New().FindUserByEmail(_jwt.UserEmail)
	exception.ResponseStatusError_New(err)

	if util.Bcrypt_CheckPasswordHash(pswd.OldPassword, user.Password) {
		user.Password, err = util.Bcrypt_HashPassword(pswd.NewPassword)
		exception.ResponseStatusError_New(err)
		err = repository.User_New().UpdateUser(user)
		exception.ResponseStatusError_New(err)

		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{})
}

func userVerified(c *gin.Context) {
	_jwt := ctr.GetJWT(c)
	if msg := _jwt.CheckingThisIsAdmin(); msg != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(msg.Error()))
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	action := c.Param("action")

	if action == "acc" {
		user, err := repository.User_New().FindUserByID(id)
		exception.ResponseStatusError_New(err)

		if !user.Verified {
			user.Verified = true
			repository.User_New().UpdateUser(user)
		}
	} else if action == "del" {
		repository.User_New().DeleteUsers(id)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func listSalesNotVerified(c *gin.Context) {
	user, err := repository.User_New().FindAllSalesNotVerified()
	exception.ResponseStatusError_New(err)

	var inInterface []map[string]interface{}
	js, err := json.Marshal(&user)
	err = json.Unmarshal(js, &inInterface)
	exception.ResponseStatusError_New(err)

	for i := 0; i < len(inInterface); i++ {
		delete(inInterface[i], "password")
	}

	c.JSON(http.StatusOK, model.HTTPResponse_Data(inInterface))
}

func registerSales(c *gin.Context) {
	var binding registerSalesJsonBinding
	err := c.ShouldBindJSON(&binding)
	exception.ResponseStatusError_New(err)

	pass, err := util.Bcrypt_HashPassword(binding.Password)
	exception.ResponseStatusError_New(err)

	user := entity.Users{
		Name:        binding.Name,
		Email:       binding.Email,
		Verified:    false,
		Password:    pass,
		PhoneNumber: binding.PhoneNumber,
		Role:        "admin",
	}
	err = repository.User_New().CreateUser(&user)
	exception.ResponseStatusError_New(err)

	c.JSON(http.StatusCreated, gin.H{})

}

func getUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := repository.User_New().FindUserByID(id)
	exception.ResponseStatusError_New(err)
	c.JSON(http.StatusOK, model.HTTPResponse_Data(user))
}

func allSales(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit == 0 {
		limit = 50
	}
	page -= 1
	if page < 0 {
		page = 0
	}
	offset := page * limit

	sales, err := repository.User_New().FindAllSales(limit, offset)
	exception.ResponseStatusError_New(err)
	totalData := repository.User_New().Total()

	var d struct {
		TotalPage int
		Sales     []model.UserSales
	}

	d.TotalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	d.Sales = sales

	c.JSON(http.StatusOK, model.HTTPResponse_Data(d))
}

func userByEmail(c *gin.Context) {
	var user selectUserBindingJSON
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	getuser, err := repository.User_New().FindUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	var pswd bool = true
	if len(getuser.Password) == 0 {
		pswd = false
	}
	c.JSON(http.StatusOK, model.HTTPResponse_Data(map[string]interface{}{
		"name":       getuser.Name,
		"phone":      getuser.PhoneNumber,
		"email":      getuser.Email,
		"registered": pswd,
	}))
}

func me(c *gin.Context) {
	_jwt := ctr.GetJWT(c)
	user, err := repository.User_New().FindUserByEmail(_jwt.UserEmail)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{})
		return
	}
	c.JSON(http.StatusOK, model.HTTPResponse_Data(map[string]interface{}{
		"name":  user.Name,
		"role":  user.Role,
		"email": user.Email,
		"id":    user.ID,
	}))
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

	fmt.Printf("%v", users)

	err = repository.User_New().CreateUser(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func updateUserSalesORAdminByAdmin(c *gin.Context) {
	var binding createUserJSONBinding
	// _jwt := ctr.GetJWT(c)

	err := c.ShouldBindJSON(&binding)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// if binding.Email != _jwt.UserEmail {
	// 	c.JSON(http.StatusBadRequest, model.HTTPResponse_Message("You not allowed update user"))
	// 	return
	// }

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

	fmt.Printf("%v", users)

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
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

	random := randStringBytesMaskImprSrc(12)

	err = service.Redis_New().SetWithExp(fmt.Sprintf("%s:changePassword", user.Email), random, REDIS_TOKEN_EXP)

	if err == nil {
		err = service.SMTP_New().SendConfirmMail("Lupa Password", user.Email, fmt.Sprintf("Password sementara anda %s, berlaku selama 5 menit", random))
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
