package main

import (
	"kala/config"
	"kala/controller"
	"kala/service"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	config.LoadEnvirolment()
	r := gin.Default()
	redis := service.Redis_New()
	r.Use(sessions.Sessions("kala_session", *redis.GetStore()))
	r.Static("/attachment", "./attachment")

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("07:00").Do(func() {
		service.GenerateAllEvidance()
	})

	controller.RegisterRoute(r)
	r.Run(":8080")
}
