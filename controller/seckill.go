package controller

import (
	"github.com/gin-gonic/gin"
	"killDemo/service"
	"strconv"
)

func GetGoodDetail(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.GetGoodDetailList(gid)
	c.JSON(res.Status, res)
}

func Initializer(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.InitializerSecKill(gid)
	c.JSON(res.Status, res)
}

func WithoutLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithoutLockSecKill(gid)
	c.JSON(res.Status, res)
}

func WithLock(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithLockSecKill(gid)
	c.JSON(res.Status, res)
}

func WithPccRead(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithPccReadSecKill(gid)
	c.JSON(res.Status, res)
}

func WithPccUpdate(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithPccUpdateSecKill(gid)
	c.JSON(res.Status, res)
}

func WithOcc(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithOccSecKill(gid)
	c.JSON(res.Status, res)
}

func WithChannel(c *gin.Context) {
	gid, _ := strconv.Atoi(c.Query("gid"))
	res := service.WithChannelSecKill(gid)
	c.JSON(res.Status, res)
}
