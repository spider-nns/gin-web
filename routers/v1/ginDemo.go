package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	//gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//日志
	ginLogFile, err := os.OpenFile("gin.log", os.O_CREATE|os.O_RDWR, 06444)
	if err != nil {
		_ = fmt.Errorf("gin log file open err,err:%v", err)
		return
	}
	//创建记录日志文件
	//gin.DefaultWriter = io.MultiWriter(ginLogFile)
	//将日志同时写入文件和控制台
	gin.DefaultWriter = io.MultiWriter(ginLogFile, os.Stdout)
	//日志格式
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s\"%s\"%s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
	simpleGroup := router.Group("/v1")
	{
		simpleGroup.GET("/get", getting)
		simpleGroup.GET("/getPath/:param", pathing)
		simpleGroup.POST("/post", posting)
		simpleGroup.POST("/mixed", mixing)
		//router.GET("/put", puting)
		//router.GET("/delete", deleting)
	}
	_ = router.Run(":8555")
}

//get+post
func mixing(context *gin.Context) {
	id := context.Query("id")
	page := context.DefaultQuery("page", "1")
	nameInForm := context.PostForm("name")
	fmt.Printf("receive param from url id:%v,page:%v,from form name:%s", id, page, nameInForm)

}

//post 参数
func posting(context *gin.Context) {
	bsData, _ := ioutil.ReadAll(context.Request.Body)
	param := &struct {
		Msg string `json:"msg"`
	}{}
	err := json.Unmarshal(bsData, param)
	if err != nil {
		log.Errorf(context, "json.Unmarshal.err, err:%v", err)
	}
	//context.Request.ParseForm()
	//msg := context.PostForm("msg")
	//设置默认值
	defaultParam := context.DefaultPostForm("name", "defaultName")
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  param.Msg,
		"data": defaultParam,
	})
}

//路径参数
func pathing(context *gin.Context) {
	nameInPath := context.Param("param")
	context.String(http.StatusOK, "hello %s from path", nameInPath)
}

//get参数
func getting(context *gin.Context) {
	name := context.DefaultQuery("name", "defaultValue")
	pathAge := context.Query("age")
	age := context.Request.URL.Query().Get("age")
	context.String(http.StatusOK, "hello %s %v %t", name, age, pathAge == age)
}
