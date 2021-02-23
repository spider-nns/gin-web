package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//禁止控制台颜色
	//1.gin.DisableConsoleColor()
	//gin.ForceConsoleColor()
	router := gin.New()

	//2.中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//3.设置表单上传大小，默认32MB
	router.MaxMultipartMemory = 1 << 20 // 1MB

	//4.日志
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
	loggerConfig := gin.LoggerConfig{
		Output:    ginLogFile,
		SkipPaths: []string{"/test"},
		Formatter: func(param gin.LogFormatterParams) string {
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
		},
	}
	router.Use(gin.LoggerWithConfig(loggerConfig))


	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})
	simpleGroup := router.Group("/v1")
	//AuthRequired just use in simpleGroup
	//simpleGroup.Use(AuthRequired())
	{
		simpleGroup.GET("/get", getting)
		simpleGroup.GET("/getPath/:param", pathing)
		simpleGroup.POST("/post", posting)
		simpleGroup.POST("/mixed", mixing)
		simpleGroup.POST("/singleFile", singleFileUpload)
		simpleGroup.POST("multiFile", multiFile)
		//router.GET("/put", puting)
		//router.GET("/delete", deleting)

		nestedGroup := simpleGroup.Group("nestedGroup")
		nestedGroup.GET("analytics", analyticsEndpoint)
		///v1/nestedGroup/analytics

	}
	_ = router.Run(":8555")
}

func analyticsEndpoint(context *gin.Context) {

}

//多文件上传
func multiFile(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		fmt.Printf("received multi file err:%v", err)
	}
	files := form.File["file[]"]
	for index, file := range files {
		log.Printf("receive no:%v file:%s", index, file.Filename)
		//
	}
	context.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	//curl -X POST  localhost:8555/v1/multiFile -H "Content-type:multipart/form-data" -F "file[]=@/Users/spider/Desktop/a.txt" -F "file[]=@/Users/spider/Desktop/年终总结.PDF"| json
}

//单文件上传
func singleFileUpload(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		fmt.Printf("received file from data err:%v", err)
	}
	fmt.Printf("receive file %s", file.Filename)
	//
	context.String(http.StatusOK, fmt.Sprintf("%s uploaded!", file.Filename))
	//curl -X POST  localhost:8555/v1/singleFile -H "Content-type:multipart/form-data" -F "file=@/Users/spider/Desktop/a.txt" | json
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
		fmt.Printf("json.Unmarshal.err, err:%v", err)
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
