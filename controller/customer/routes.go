package customer

import (
	"kala/model"
	"kala/repository"
	"kala/repository/entity"
	"net/http"
	"strconv"
	"time"

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
}

type jsonUpdateCustomer struct {
	ID            int       `json:"ID"`
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
}

func ResgisterRoutes(c *gin.RouterGroup, ctx *gin.Engine) {
	c.POST("/customer", createCustomer)
	c.DELETE("/customer", deleteCustomer)
	c.GET("/customer/sales/:id", findBySalesID)
	c.GET("/customer/:id", findByCustomerID)
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

	err = repository.CustomerRepository_New().DeleteCustomer(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func findBySalesID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cus, err := repository.CustomerRepository_New().FindCustomerBySalesID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.HTTPResponse_Data(cus))
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
