package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	docs "test-app/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const PORT = 8080

//	@Summary		Hello
//	@Router			/ [get]
//	@x-perf-check	{ "latency": 100, "errorRate": 0.1 }
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

//	@Summary		Sleep for n seconds
//	@Router			/sleep [get]
//	@x-perf-check	{ "latency": 250, "errorRate": 0.2, "params": { "query": { "seconds": { "examples": [1] } } } }
func SleepPath(g *gin.Context) {
	secondsString := g.Query("seconds")

	seconds, err := strconv.Atoi(secondsString)
	if err != nil {
		panic("seconds must be an integer")
	}

	time.Sleep(time.Second * time.Duration(seconds))
	g.JSON(http.StatusOK, "sleep-route")
}

//	@Summary		Example param
//	@Router			/param/{x} [get]
//	@Param			x	path	int	true	"X param"
//	@x-perf-check	{  "latency": 100, "errorRate": 0.1, "params": { "path": { "x": {"examples": ["abc", "def"] } } } }
func ExampleParamPath(g *gin.Context) {
	x := g.Param("x")

	g.JSON(http.StatusOK, fmt.Sprintf("param is %s", x))
}

//	@Summary		Range param
//	@Router			/range/{x} [get]
//	@Param			x	path	int	true	"X param"
//	@x-perf-check	{  "latency": 100, "errorRate": 0.1, "params": { "path": { "x": {"range": { "min": 0, "max": 1000 } } } } }
func RangeParamPath(g *gin.Context) {
	x := g.Param("x")

	g.JSON(http.StatusOK, fmt.Sprintf("param is %s", x))
}

//	@Summary		Pattern param
//	@Router			/pattern/{x}/{y} [get]
//	@Param			x	path	int	true	"X param"
//	@x-perf-check	{  "latency": 100, "errorRate": 0.1, "params": { "path": { "x": {"pattern": "uuid" }, "y": { "pattern": "string(8)" } } } }
func PatternParamPath(g *gin.Context) {
	x := g.Param("x")
	y := g.Param("y")

	g.JSON(http.StatusOK, fmt.Sprintf("params are %s, %s", x, y))
}

//	@title			Example API
//	@schemes		http
//	@x-perf-check	{ "users": 100, "duration": 3 }
func main() {
	r := gin.Default()

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", PORT)

	r.GET("/", Helloworld)

	r.GET("/sleep", SleepPath)

	r.GET("/param/:x", ExampleParamPath)

	r.GET("/range/:x", RangeParamPath)

	r.GET("/pattern/:x/:y", PatternParamPath)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(fmt.Sprintf(":%d", PORT))
}
