package service

import (
	"errors"
	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"killDemo/cache"
	"killDemo/dao"
	"killDemo/tool"
	"killDemo/tool/e"
	"strconv"
	"time"
)

func WithRedisList(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	g := strconv.Itoa(gid)
	for i := 0; i < seckillNum; i++ {
		cache.RedisClient.LPush(g, g)
	}
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithRedisListSecKillGoods(gid, userID)
			if err != nil {
				code = e.ERROR
				logging.Errorln("Error", err)
			} else {
				logging.Infof("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	err := AfterRedisListSecKill(gid)
	kCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Infoln("Error")
	}
	logging.Infof("Total %v goods", kCount)
	return tool.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithRedisListSecKillGoods(gid, userID int) error {
	g := strconv.Itoa(gid)
	u := strconv.Itoa(userID)
	if cache.RedisClient.Get(u+g).Val() == "" { // 这用户没有秒杀过
		cache.RedisClient.RPop(g) // 表示库存 -1
		cache.RedisClient.Set(u+g, g, 2*time.Minute)
		cache.RedisClient.ZAdd(g, redis.Z{float64(time.Now().Unix()), userID})
	} else { // 这用户已经有记录了
		return errors.New("该用户已经抢过了")
	}
	return nil
}

func AfterRedisListSecKill(gid int) error {
	g := strconv.Itoa(gid)
	// ZRevRangeWithScores用于获取有序集合（sorted set）g中分数从高到低排列的所有成员及其对应的分数。
	// g是有序集合的键名，0表示从分数最高的成员开始获取，-1表示获取到分数最低的成员。
	// 即这个命令获取了整个有序集合的所有成员及其分数。
	ret, _ := cache.RedisClient.ZRevRangeWithScores(g, 0, -1).Result()
	for _, z := range ret {
		userID, err := strconv.Atoi(z.Member.(string))
		tx := dao.DB.Begin()
		// 1. 扣库存
		err = dao.ReduceOneByGoodsId(gid)
		if err != nil {
			tx.Rollback()
			return err
		}
		// 2. 创建订单
		kill := dao.SuccessKilled{
			GoodsId:    int64(gid),
			UserId:     int64(userID),
			State:      0,
			CreateTime: time.Now(),
		}
		err = dao.CreateOrder(kill)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}
