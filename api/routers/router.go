package routers

import (
	"remix-api/api/http/controllers"
	"remix-api/api/http/middlewares"
	"remix-api/configs"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	logger := new(middlewares.Logger)
	allowCors := new(middlewares.AllowCors)
	authCheck := new(middlewares.AuthCheck)
	router.Use(logger.Handle())
	router.Use(gin.Recovery())
	router.Use(allowCors.Handle())
	throttle := middlewares.NewThrottle(configs.Config.MaxThrottle)
	gin.SetMode(configs.Config.Env)

	api := router.Group("/api/v2")
	noAuthApi := router.Group("/api/v2")
	api.Use(authCheck.Handle()).Use(throttle.Handle())
	fileApi := router.Group("/api/v2")

	noPrefix := fileApi.Group("")
	{
		controller := new(controllers.Controller)
		noPrefix.GET("files", controller.GetFiles)
	}

	task := api.Group("task")
	{
		taskController := new(controllers.TaskController)
		// 创建任务
		task.POST("", taskController.Create)
		// 任务列表
		task.GET("", taskController.GetList)
		// 查询任务
		task.GET("show", taskController.Show)
		// 下载任务文件
		task.GET("download", taskController.Download)
		// 删除任务
		task.DELETE("", taskController.Delete)
	}

	noAuth := noAuthApi.Group("")
	{
		userController := new(controllers.UserController)
		configController := new(controllers.ConfigController)
		noAuth.POST("/mini-program-login", userController.MiniProgramLogin)
		noAuth.GET("/config/:id", configController.Show)
	}
	return router
}
