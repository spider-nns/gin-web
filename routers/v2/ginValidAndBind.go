package main

import (
	"gin-web/routers/v2/ao"
	"gin-web/routers/v2/custValidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
)

//gin valid and bind
//MustBindWith，如果存在绑定错误，请求将被以下指令中止 c.AbortWithError(400, err).SetType(ErrorTypeBind)
//响应状态代码会被设置为400，请求头Content-Type被设置为text/plain; charset=utf-8
//注意，如果你试图在此之后设置响应代码，将会发出一个警告
//[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 422
//如果你希望更好地控制行为，请使用ShouldBind相关的方法

//当我们使用绑定方法时，Gin会根据Content-Type推断出使用哪种绑定器，如果你确定你绑定的是什么，你可以使用MustBindWith或者BindingWith。
//你还可以给字段指定特定规则的修饰符，如果一个字段用binding:"required"修饰，并且在绑定时该字段的值为空，那么将返回一个错误
func main() {
	//1.new
	r := gin.New()

	//2.中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("BookableDate", custValidator.BookableDate)
	}
	//router
	v2 := r.Group("/v2")
	{
		v2.POST("/loginJson", validJson)
		// curl -X POST localhost:8556/loginJson -d '{"userName":"guest","password":"guest"}' | json
		v2.POST("/loginForm", validForm)
		v2.GET("/customerValid", customerValidHandler)
	}
	_ = r.Run(":8556")
}

//自定义验证
func customerValidHandler(context *gin.Context) {
	var b custValidator.Booking
	if err := context.ShouldBindWith(&b, binding.Query); err == nil {
		context.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func validJson(context *gin.Context) {
	var json ao.Login
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.UserName != "guest" || json.Password != "guest" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "welcome,sir"})
}
func validForm(context *gin.Context) {
	var form ao.Login
	if err := context.ShouldBind(&form); err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if form.UserName != "guest" || form.Password != "guest" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "welcome,sir"})
}
