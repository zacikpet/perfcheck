package main

import (
	"net/http"

	docs "test-app/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@Summary		Hello
//	@Router			/ [get]
//	@x-perf-check	{ "latency": 100, "errorRate": 0.1 }
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

//	@Summary		TestRoute
//	@Router			/test-route [get]
//	@x-perf-check	{ "latency": 250, "errorRate": 0.2 }
func TestRoute(g *gin.Context) {
	g.JSON(http.StatusOK, "test-route")
}

//	@Summary		TestRoute
//	@Router			/group/a [post]
//	@x-perf-check	{ "latency": 150, "errorRate": 0.05 }
func GroupA(g *gin.Context) {
	g.JSON(http.StatusOK, "group/a")
}

//	@Summary		TestRoute
//	@Router			/group/b [patch]
//	@x-perf-check	{ "latency": 150, "errorRate": 0.1 }
func GroupB(g *gin.Context) {
	g.JSON(http.StatusOK, "group/b")
}

func main() {
	r := gin.Default()

	docs.SwaggerInfo.Title = "Swagger Example API"

	r.GET("/", Helloworld)

	r.GET("/test-route", TestRoute)

	group := r.Group("group")
	{
		group.POST("/a", GroupA)
		group.PATCH("/b", GroupB)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(":8080")
}
