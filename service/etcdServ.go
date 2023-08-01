package service

import (
	"context"
	"fmt"
	logging "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"killDemo/tool"
	"killDemo/tool/e"
	"time"
)

var (
	defaultEtcdConfig = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
)

type EtcdMutex struct {
	Ttl     int64              //租约时间
	Conf    clientv3.Config    //etcd集群配置
	Key     string             //etcd的key
	cancel  context.CancelFunc //关闭续租的func
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	txn     clientv3.Txn
}

func (em *EtcdMutex) initETCD() error {
	var err error
	var ctx context.Context
	client, err := clientv3.New(em.Conf)
	if err != nil {
		fmt.Println("New etcd error ", err)
		return err
	}
	em.txn = clientv3.NewKV(client).Txn(context.TODO())
	em.lease = clientv3.NewLease(client)
	leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
	if err != nil {
		return err
	}

	ctx, em.cancel = context.WithCancel(context.TODO())
	em.leaseID = leaseResp.ID
	_, err = em.lease.KeepAlive(ctx, em.leaseID)
	return err
}

func (em *EtcdMutex) Lock() error {
	err := em.initETCD()
	if err != nil {
		return err
	}

	em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)).
		Then(clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseID))).
		Else()
	txnResp, err := em.txn.Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded { //判断txn.if条件是否成立
		return fmt.Errorf("抢锁失败")
	}

	return nil
}

func (em *EtcdMutex) UnLock() {
	em.cancel()
	_, _ = em.lease.Revoke(context.TODO(), em.leaseID)
	fmt.Println("释放了锁")
}

func WithETCDSecKillGoods(gid, userID int) error {

	eMutex1 := &EtcdMutex{
		Conf: defaultEtcdConfig,
		Ttl:  10,
		Key:  "lock",
	}
	err := eMutex1.Lock()
	if err != nil {
		return err
	}
	err = WithoutLockSecKillGoods(gid, userID)
	eMutex1.UnLock()
	return err
}

func WithETCDSecKill(gid int) tool.Response {
	code := e.SUCCESS
	seckillNum := 50
	wg.Add(seckillNum)
	InitializerSecKill(gid)

	for i := 0; i < seckillNum; i++ {
		userID := i
		go func() {
			err := WithETCDSecKillGoods(gid, userID)
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
