package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	//自定义中间件
	r.Use(Logger())
	//http 重定向
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://google.com")
	})
	//路由重定向
	r.GET("/test1", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})
	r.GET("/middleWare", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		log.Println(example)
	})
	_ = r.Run(":8558")
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t:= time.Now()

		c.Set("example","123456")

		//
		c.Next()

		//after request

		latency := time.Since(t)
		log.Print(latency)

		//access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
