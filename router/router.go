package router

import (
	"github.com/gin-gonic/gin"
	"killDemo/controller"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 商品信息展示页面获取数据
	r.GET("/good", controller.GetGoodDetail)
	r.GET("/creat", controller.Initializer)

	skGroup := r.Group("/v1")
	{
		// 不加锁,出现超卖现象
		skGroup.GET("/without-lock", controller.WithoutLock)

		// 加锁(sync包中的Mutex类型的互斥锁),没有问题
		skGroup.GET("/with-lock", controller.WithLock)
		// 加锁(数据库悲观锁，查询加锁),超卖
		skGroup.GET("/with-pcc-read", controller.WithPccRead)
		// 加锁(数据库悲观锁，更新限定), 正常
		skGroup.GET("/with-pcc-update", controller.WithPccUpdate)
		// 加锁(数据库乐观锁，正常)
		skGroup.GET("/with-occ", controller.WithOcc)
		// channel 限制，正常
		skGroup.GET("/with-channel", controller.WithChannel)
	}
	// 分布式
	skDisGroup := r.Group("/v2")
	{
		skDisGroup.GET("/rush", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		// 基于redis的redission分布式,正常
		skDisGroup.GET("/with-redission", controller.WithRedission)
		// 基于ETCD的锁, 正常
		skDisGroup.GET("/with-etcd", controller.WithETCD)
		// 基于redis的List, 正常
		skDisGroup.GET("/with-redis-list", controller.WithRedisList)
	}
	return r
}
