package customer

import (
	"fmt"
	"kala/exception"
	"kala/model"
	"kala/repository"
	"kala/repository/entity"
	"math"
	"net/http"
	"strconv"
	"time"

	ctr "kala/controller/util"

	"github.com/gin-gonic/gin"
)

type jsonAddCustomer struct {
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	NoHp          string    `json:"no_hp"`
	TglDec        time.Time `json:"tgl_dec"`
	TglLahir      time.Time `json:"tgl_lahir"`
	TglSTNK       time.Time `json:"tgl_stnk"`
	TglAngsuran   time.Time `json:"tgl_angsuran"`
	NoRangka      string    `json:"no_rangka"`
	TypeKendaraan string    `json:"type_kendaraan"`
	Leasing       string    `json:"leasing"`
	SalesID       int       `json:"sales_id"`
	TypeAngsuran  string    `json:"type_angsuran"`
	TotalAngsuran int       `json:"total_angsuran"`
}

type jsonUpdateCustomer struct {
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	NoHp          string    `json:"no_hp"`
	TglDec        time.Time `json:"tgl_dec"`
	TglLahir      time.Time `json:"tgl_lahir"`
	TglSTNK       time.Time `json:"tgl_stnk"`
	TglAngsuran   time.Time `json:"tgl_angsuran"`
	NoRangka      string    `json:"no_rangka"`
	TypeKendaraan string    `json:"type_kendaraan"`
	Leasing       string    `json:"leasing"`
	SalesID       int       `json:"sales_id"`
	TypeAngsuran  string    `json:"type_angsuran"`
	TotalAngsuran int       `json:"total_angsuran"`
}

func ResgisterRoutes(c *gin.RouterGroup, ctx *gin.Engine) {
	c.POST("/customer", createCustomer)
	c.DELETE("/customer/:id", deleteCustomer)
	c.GET("/customer/sales", findBySalesID)
	c.GET("/customer/:id", findByCustomerID)
	c.GET("/customer", listCustomer)
	c.PATCH("/customer/:id", editCustomer)
	// "/customer/all-sales",
}

func editCustomer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param(("id")))
	var data jsonUpdateCustomer

	err := c.ShouldBindJSON(&data)
	exception.ResponseStatusError_New(err)
	cus, err := repository.CustomerRepository_New().FindCustomerByID(id)
	exception.ResponseStatusError_New(err)

	customer := cus.Customers

	customer.Address = data.Address
	customer.Leasing = data.Leasing
	customer.Name = data.Name
	customer.NoHp = data.NoHp
	customer.TglAngsuran = data.TglAngsuran
	customer.TglDec = data.TglDec
	customer.TglSTNK = data.TglSTNK
	customer.TglLahir = data.TglLahir
	customer.TotalAngsuran = data.TotalAngsuran
	customer.TypeAngsuran = data.TypeAngsuran
	customer.TypeKendaraan = data.TypeKendaraan
	customer.SalesID = data.SalesID
	// customer.ID = id
	fmt.Printf("%+v", customer)
	err = repository.CustomerRepository_New().UpdateCustomer(&customer)
	exception.ResponseStatusError_New(err)
	c.JSON(http.StatusOK, gin.H{})
}

func listCustomer(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "50"))

	if err != nil {
		limit = 50
	}
	page -= 1
	customers, err := repository.CustomerRepository_New().ListAllCustomer(page*limit, limit)
	exception.ResponseStatusError_New(err)
	totalData := repository.CustomerRepository_New().Total()

	var d struct {
		TotalPage int
		Customer  []entity.CustomerInnerJoinUser
	}

	d.TotalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	d.Customer = customers

	ctx.JSON(http.StatusOK, model.HTTPResponse_Data(d))
}

func createCustomer(c *gin.Context) {
	var customer jsonAddCustomer

	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cus := entity.Customers{
		ID:            0,
		Name:          customer.Name,
		Address:       customer.Address,
		NoHp:          customer.NoHp,
		TglDec:        customer.TglDec,
		TglLahir:      customer.TglLahir,
		TglAngsuran:   customer.TglAngsuran,
		TglSTNK:       customer.TglSTNK,
		NoRangka:      customer.NoRangka,
		TypeKendaraan: customer.TypeKendaraan,
		Leasing:       customer.Leasing,
		SalesID:       customer.SalesID,
		NewCustomer:   true,
		TypeAngsuran:  customer.TypeAngsuran,
		TotalAngsuran: customer.TotalAngsuran,
	}
	err = repository.CustomerRepository_New().CreateCustomer(&cus)

	if err != nil {
		c.JSON(http.StatusBadGateway, model.HTTPResponse_Message(err.Error()))
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func deleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository.Notification_New().DeleteByCustomerID(id)
	repository.EvidanceRepository_New().DeleteByCustomerID(id)

	err = repository.CustomerRepository_New().DeleteCustomer(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func findBySalesID(c *gin.Context) {
	_jwt := ctr.GetJWT(c)

	user, err := repository.User_New().FindUserByEmail(_jwt.UserEmail)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

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

	cus, err := repository.CustomerRepository_New().FindCustomerBySalesID(offset, limit, user.ID)
	totalData := repository.CustomerRepository_New().TotalCustomerBySalesID(user.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var d struct {
		TotalPage int
		Customer  []entity.Customers
	}

	d.TotalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	d.Customer = cus

	c.JSON(http.StatusOK, model.HTTPResponse_Data(d))
}

func findByCustomerID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cus, err := repository.CustomerRepository_New().FindCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.HTTPResponse_Data(cus))
}
