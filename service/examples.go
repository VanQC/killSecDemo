package service

import (
	"fmt"
	logging "github.com/sirupsen/logrus"
	"killDemo/tool"
	"killDemo/tool/e"
)

func WithoutLockSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithoutLockSecKillGoods(gid, userID)
			if err != nil {
				fmt.Println("Error", err)
			} else {
				fmt.Printf("User: %d seckill successfully.\n", userID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	killedCount, err := GetKilledCount(gid)
	if err != nil {
		code = e.ERROR
		logging.Error("Seckill System Error")
		return tool.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	fmt.Println(killedCount)
	logging.Infof("kill %v product", killedCount)
	return tool.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func WithLockSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithLockSecKillGoods(gid, userID)
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

func WithPccReadSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithPccReadSecKillGoods(gid, userID)
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

func WithPccUpdateSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithPccUpdateSecKillGoods(gid, userID)
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

func WithOccSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithOccSecKillGoods(gid, userID, 1)
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

func WithChannelSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 200
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	go ChannelConsumer()
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithChannelSecKillGoods(gid, userID)
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
