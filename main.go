package main

import (
	"fmt"
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
		// s.Every(90).Seconds().Do(func() {
		fmt.Println("create evidance")
		service.GenerateAllEvidance()
	})
	s.StartAsync()

	controller.RegisterRoute(r)
	r.Run(":8080")
	fmt.Println("Stoped")
	s.Stop()
}
