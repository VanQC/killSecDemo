package service

import (
	"errors"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"killDemo/dao"
	"killDemo/tool"
	"killDemo/tool/e"
	"sync"
	"time"
)

// 获取商品的详细信息
func GetGoodDetailList(gid int) tool.Response {
	code := e.SUCCESS
	good, err := dao.FindGoodsById(gid)
	if err != nil {
		code = e.ERROR
		return tool.Response{
			Status: code,
			Error:  err.Error(),
		}
	}
	return tool.Response{
		Status: code,
		Data:   good,
		Msg:    e.GetMsg(code),
	}
}

var wg sync.WaitGroup
var lock sync.Mutex

func InitializerSecKill(gid int) tool.Response {
	code := e.SUCCESS

	tx := dao.DB.Begin()            // 开启事务
	err := dao.DeleteByGoodsId(gid) // 删除事务
	if err != nil {                 // 发生错误的话就进行回滚
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return tool.Response{
				Status: 500,
				Msg:    "删除发生错误，且不是 RecordNotFound",
				Error:  err.Error(),
			}
		}
	}
	err = dao.UpdateCountByGoodsId(gid) // 更新事务
	if err != nil {
		tx.Rollback()
		return tool.Response{
			Status: 500,
			Msg:    "更新事务出错",
			Error:  err.Error(),
		}
	}
	tx.Commit()
	return tool.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 获取总共秒杀了多少商品
func GetKilledCount(gid int) (int64, error) {
	return dao.GetKilledCountByGoodsId(gid)
}

func WithoutLockSecKillGoods(gid, userID int) error {
	tx := dao.DB.Begin()
	// 检查库存
	count, err := dao.SelectCountByGoodsId(gid)
	if err != nil {
		return err
	}
	if count > 0 {
		// 1. 扣库存
		err = dao.ReduceStockByGoodsId(gid, int(count-1))
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
	}
	tx.Commit()
	return nil
}

func WithLockSecKillGoods(gid, userID int) error {
	lock.Lock()
	err := WithoutLockSecKillGoods(gid, userID)
	lock.Unlock()
	return err
}

func WithPccReadSecKillGoods(gid, userID int) error {
	tx := dao.DB.Begin()
	count, err := dao.SelectCountByGoodsIdPcc(gid)
	// 先读后更新的数据竞争场景
	if err != nil {
		return err
	}
	if count > 0 {
		// 1. 扣库存
		err = dao.ReduceStockByGoodsId(gid, int(count-1))
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
	}
	tx.Commit()
	return nil
}

func WithPccUpdateSecKillGoods(gid, userID int) error {
	tx := dao.DB.Begin()
	// 1. 扣库存
	count, err := dao.ReduceByGoodsId(gid)
	if err != nil {
		return err
	}
	if count > 0 {
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
	}
	tx.Commit()
	return nil
}

func WithOccSecKillGoods(gid, userID, num int) error {
	tx := dao.DB.Begin()
	good, err := dao.SelectGoodByGoodsId(gid)
	if err != nil {
		return err
	}
	if int(good.PsCount) >= num {
		// 1. 扣库存
		count, err := dao.ReduceStockByOcc(gid, num, int(good.Version))
		if err != nil {
			tx.Rollback()
			return err
		}
		if count > 0 {
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
		} else {
			tx.Rollback()
		}
	} else {
		tx.Rollback()
		return errors.New("库存不够了")
	}
	tx.Commit()
	return nil
}

func WithChannelSecKillGoods(gid, userID int) error {
	kill := [2]int{gid, userID}
	kChan := GetInstance()
	*kChan <- kill
	return nil
}

func ChannelConsumer() {
	for {
		kill, ok := <-(*GetInstance())
		if !ok {
			continue
		}
		err := WithoutLockSecKillGoods(kill[0], kill[1])
		if err != nil {
			logging.Error("Error")
		} else {
			logging.Infof("User:%v SecKill Successfully", kill[1])
		}
	}
}

type SingleTon chan [2]int

var instance *SingleTon
var once sync.Once

func GetInstance() *SingleTon {
	once.Do(func() {
		ret := make(SingleTon, 100)
		instance = &ret
	})
	return instance
}
