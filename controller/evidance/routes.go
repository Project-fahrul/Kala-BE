package evidance

import (
	"fmt"
	"kala/exception"
	"kala/model"
	"kala/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(c *gin.RouterGroup) {
	c.POST("/evidance", uploadFile)
	c.GET("/evidance", listEvidance)
	c.GET("/evidance/:type", getEvidance)
	c.GET("evidance/count", func(ctx *gin.Context) {
		c, e := repository.EvidanceRepository_New().Count()
		exception.ResponseStatusError_New(e)

		ctx.JSON(http.StatusOK, model.HTTPResponse_Data(c))
	})
}

func getEvidance(c *gin.Context) {
	typeEvidance := c.Param("type")
	sales_id, _ := strconv.Atoi(c.Query("sales"))
	customer_id, _ := strconv.Atoi(c.Query("customer"))
	due, err := time.Parse("2006-01-02", c.Query("due"))

	exception.ResponseStatusError_New(err)
	e, err := repository.EvidanceRepository_New().Evidance(sales_id, customer_id, due, typeEvidance)
	exception.ResponseStatusError_New(err)

	c.JSON(http.StatusOK, model.HTTPResponse_Data(e))
}

func uploadFile(c *gin.Context) {
	customer_id := c.PostForm("customer_id")
	sales_id := c.PostForm("sales_id")
	message := c.PostForm("message")
	typeEvidance := c.PostForm("type")
	image, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	fmt.Printf("%s,%s,%s,%s,%s\n", customer_id, sales_id, message, typeEvidance, image.Filename)
	c.JSON(http.StatusCreated, gin.H{})
}

func listEvidance(c *gin.Context) {
	ev, err := repository.EvidanceRepository_New().ListEvidance()

	fmt.Printf("%v", ev)
	exception.ResponseStatusError_New(err)
	c.JSON(http.StatusOK, model.HTTPResponse_Data(ev))
}
