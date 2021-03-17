package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/controller"
	"net/http"
	"strconv"
)

func router(publicPath string) http.Handler {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile(publicPath, true)))

	list := router.Group("/api/list")

	list.GET("/all", handlerAll)
	list.GET("/top10", handlerTop10)
	list.GET("/bottom10", handlerBottom10)
	list.GET("/inRank", handlerInRank)
	list.GET("/getBySymbol", handlerGetBySymbol)

	return router

}

func response(c *gin.Context, data interface{}, err error) {
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func handlerGetBySymbol(c *gin.Context) {
	symbols, _ := c.GetQueryArray("symbols")

	list, err := controller.Instance().GetBySymbol(symbols)
	response(c, list, err)
}

func handlerInRank(c *gin.Context) {
	start, _ := c.GetQuery("start")
	end, _ := c.GetQuery("end")

	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)

	list, err := controller.Instance().InRank(startInt, endInt)
	response(c, list, err)
}

func handlerBottom10(c *gin.Context) {
	list, err := controller.Instance().Bottom10()
	response(c, list, err)
}
func handlerTop10(c *gin.Context) {
	list, err := controller.Instance().Top10()
	response(c, list, err)
}

func handlerAll(c *gin.Context) {
	list, err := controller.Instance().All()
	response(c, list, err)
}
