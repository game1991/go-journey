package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "p", "9990", "端口")
}

func main() {
	flag.Parse()
	if port == "" {
		panic("-p传参对象不能为空")
	}

	engine := gin.Default()

	route := engine.Group("")

	route.GET("", func(c *gin.Context) {

	})
	engine.Run(":" + port)
}
