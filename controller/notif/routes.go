package notif

import (
	ctr "kala/controller/util"
	"kala/model"
	"kala/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(c *gin.RouterGroup) {
	c.GET("/notif", ListAllNotifBySalesID)
}

func ListAllNotifBySalesID(c *gin.Context) {
	_jwt := ctr.GetJWT(c)

	user, err := repository.User_New().FindUserByEmail(_jwt.UserEmail)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPResponse_Message(err.Error()))
		return
	}

	notif, err := repository.Notification_New().ListAllNotificationBySalesID(user.ID)

	c.JSON(http.StatusOK, model.HTTPResponse_Data(notif))
}
