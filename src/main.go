package main

import (
	"net/http"

	"github.com//controller"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	r := gin.Default()
// 	r.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello World")
// 	})
// 	r.Run() // listen and server on 0.0.0.0:8080
// }

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "aaa")
	})
	router.GET("/tasks", controller.TasksGET)
	// v := router.Group("/v1")
	// {
	// 	v.GET("/", func(c *gin.Context) {
	// 		c.String(http.StatusOK, "aaa")
	// 	})
	// 	v.GET("/tasks", controller.TasksGET)
	// 	v.POST("/tasks", controller.TaskPOST)
	// 	v.PATCH("/tasks/:id", controller.TaskPATCH)
	// 	v.DELETE("/tasks/:id", controller.TaskDELETE)
	// }
	// v1 := router.Group("/task-checker/api/v1")
	// {
	// 	v1.GET("/", func(c *gin.Context) {
	// 		c.String(http.StatusOK, "aaa")
	// 	})
	// 	v1.GET("/tasks", controller.TasksGET)
	// 	v1.POST("/tasks", controller.TaskPOST)
	// 	v1.PATCH("/tasks/:id", controller.TaskPATCH)
	// 	v1.DELETE("/tasks/:id", controller.TaskDELETE)
	// }
	// nginxのreverse proxy設定
	router.Run()
}
