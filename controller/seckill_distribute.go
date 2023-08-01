package controller

import (
	"github.com/gin-gonic/gin"
	"killDemo/service"
	"strconv"
)

func WithRedission(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithRedissionSecKill(gid)
	c.JSON(res.Status, res)
}

func WithRedisList(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithRedisList(gid)
	c.JSON(res.Status, res)
}

func WithETCD(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithETCDSecKill(gid)
	c.JSON(res.Status, res)
}
