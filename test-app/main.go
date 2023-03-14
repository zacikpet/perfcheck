package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	docs "test-app/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const PORT = 8080

// @Summary		Hello
// @Router			/ [get]
// @x-perf-check	{ "latency": ["avg < 50", "min < 50", "avg_stat < 50"], "errorRate": ["avg_stat < 0.1"] }
func Helloworld(g *gin.Context) {

	sleep := rand.Intn(100)

	time.Sleep(time.Millisecond * time.Duration(sleep))
	g.JSON(http.StatusOK, "helloworld")
}

// @title			Example API
// @schemes		http
// @x-perf-check	{ "stages": [{ "duration": "1s", "target": 5 }] }
func main() {
	r := gin.Default()

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", PORT)

	r.GET("/", Helloworld)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(fmt.Sprintf(":%d", PORT))
}
