package evidance

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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	customer_id, _ := strconv.Atoi(c.PostForm("customer_id"))
	sales_id, _ := strconv.Atoi(c.PostForm("sales_id"))
	message := c.PostForm("message")
	typeEvidance := c.PostForm("type")
	image, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	sales, err := repository.User_New().FindUserByID(sales_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	id := uuid.New()
	fileName := fmt.Sprintf("%s-%s-%s", sales.Name, id.String(), image.Filename)

	err = c.SaveUploadedFile(image, "./attachment/"+fileName)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	exception.ResponseStatusError_New(err)
	ev := entity.Evidances{
		SalesID:      sales_id,
		CustomerID:   customer_id,
		SubmitDate:   time.Now(),
		TypeEvidance: typeEvidance,
		Comment:      message,
		Content:      fileName,
	}

	// repository.Notification_New().Delete(ev)

	err = repository.EvidanceRepository_New().UploadFile(ev)
	exception.ResponseStatusError_New(err)
	c.JSON(http.StatusCreated, gin.H{})
}

func listEvidance(c *gin.Context) {

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
	totalData := repository.EvidanceRepository_New().Total()

	ev, err := repository.EvidanceRepository_New().ListEvidanceWithLimit(limit, offset)

	var d struct {
		TotalPage int
		Evidance  []model.ListEvidance
	}

	d.TotalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	d.Evidance = ev

	exception.ResponseStatusError_New(err)
	c.JSON(http.StatusOK, model.HTTPResponse_Data(d))
}
