package service

import (
	"bytes"
	"errors"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"killDemo/cache"
	"killDemo/tool"
	"killDemo/tool/e"
	"math/rand"
	"strconv"
	"time"
)

func WithRedissionSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)
	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithRedssionSecKillGoods(gid, userID)
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

func WithRedssionSecKillGoods(gid, userID int) error {
	g := strconv.Itoa(gid)
	verifyValue := getUuid(g)
	for {
		lockSuccess, err := cache.RedisClient.SetNX(g, verifyValue, time.Second*3).Result()
		if err != nil {
			return errors.New("get lock fail")
		}
		if lockSuccess {
			break
		}
	}
	fmt.Println("get lock success")

	err := WithoutLockSecKillGoods(gid, userID)
	if err != nil {
		return err
	}
	value, _ := cache.RedisClient.Get(g).Result()
	if value == verifyValue { //compare value,if equal then del
		_, err := cache.RedisClient.Del(g).Result()
		if err != nil {
			fmt.Println("unlock fail")
			return nil
		} else {
			fmt.Println("unlock success")
		}
	}
	return nil
}

func getUuid(gid string) string {
	codeLen := 8
	// 1. 定义原始字符串
	rawStr := "jkwangagDGFHGSERKILMJHSNOPQR546413890_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := r.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String() + gid
}
