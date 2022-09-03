package main

import (
	"kala/config"
	"kala/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnvirolment()
	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("kala_session", store))

	controller.RegisterRoute(r)
}
