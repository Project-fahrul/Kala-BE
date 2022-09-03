package util

import (
	"kala/util"

	"github.com/gin-gonic/gin"
)

func GetJWT(c *gin.Context) *util.JWT {
	val, ok := c.Get("auth")
	if !ok {
		return nil
	}
	return val.(*util.JWT)
}
