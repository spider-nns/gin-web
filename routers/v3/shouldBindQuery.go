package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	//设置静态数据
	r.Static("/assets", "./assets")
	r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	bind := r.Group("/bind")
	bind.Any("/sbq",shouldBindQuery)
	bind.GET("/json",jsonReturn)
	bind.GET("/moreJson",moreJson)

	bind.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	bind.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	bind.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		// The specific definition of protobuf is written in the testdata/protoexample file.
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		// Note that data becomes binary data in the response
		// Will output protoexample.Test protobuf serialized data
		c.ProtoBuf(http.StatusOK, data)
	})

	r.Run(":8557")
}

func moreJson(c *gin.Context) {
	r := new(Result)
	r.Code=200
	r.Msg="200"
	r.Data=""
	c.JSON(http.StatusOK,r)
}

func jsonReturn(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message": "hey", "status": http.StatusOK})
}

type Person struct {
	Name string `form:"name" json:"name"`
	Pass string `form:"pass" json:"pass"`
}
type Result struct {
	Code uint8 `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
func shouldBindQuery(c *gin.Context) {
	var p Person
	if err := c.ShouldBindQuery(&p);err == nil{
		log.Println("====== Only Bind By Query String ======")
		log.Println(p.Name)
		log.Println(p.Pass)
	}
	if err := c.Bind(&p);err == nil{
		log.Println("====== Bind ======")
		log.Println(p.Name)
		log.Println(p.Pass)
	}
	if err := c.BindJSON(&p);err == nil {
		log.Println("====== BindJson ======")
		log.Println(p.Name)
		log.Println(p.Pass)
	}
	if err := c.ShouldBind(&p);err == nil {
		log.Println("====== ShouldBind ======")
		log.Println(p.Name)
		log.Println(p.Pass)
	}
	c.String(http.StatusOK,"success")
}
