package main

import (
	"kala/config"
	"kala/controller"
	"kala/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnvirolment()
	r := gin.Default()
	redis := service.Redis_New()
	r.Use(sessions.Sessions("kala_session", *redis.GetStore()))

	controller.RegisterRoute(r)
	r.Run(":8080")
}
