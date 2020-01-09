package zkLock

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"zklock/lock"
	"zklock/util"

	"github.com/samuel/go-zookeeper/zk"
)

const EPHEMERAL_SEQUENTIAL = 3

type ZkLock struct {
	Paths    []string
	LockPath string
	c        *zk.Conn
}

func NewLock(opt ...lock.Option) (*ZkLock, error) {
	optConf := &lock.Options{}
	for _, o := range opt {
		o(optConf)
	}
	lock := &ZkLock{
		Paths:    optConf.Paths,
		LockPath: optConf.LockPath,
	}
	if lock.LockPath == "" {
		lock.LockPath = "/my_lock"
	}

	conn, _, err := zk.Connect(lock.Paths, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("连接错误%v", err)
	}
	lock.c = conn
	return lock, nil
}

func (n *ZkLock) TryLock() (bool, error) {
	fmt.Println("开始尝试获取分布式锁--------------")
	exist, _, err := n.c.Exists(fmt.Sprintf("/%s", n.LockPath))
	if err != nil {
		return false, fmt.Errorf("检查节点是否存在失败,err:%v", err)
	}
	if !exist {
		n.c.Create(fmt.Sprintf("/%s", n.LockPath), []byte("locked"), 0, zk.WorldACL(zk.PermAll))
	}
	// 创建当前锁节点
	currentNode, err := n.c.Create(fmt.Sprintf("/%s/locked", n.LockPath), []byte("locked"), EPHEMERAL_SEQUENTIAL, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, fmt.Errorf("创建节点失败,err:%v", err)
	}
	//获取所有子节点
	childrens, _, err := n.c.Children(fmt.Sprintf("/%s", n.LockPath))
	if err != nil {
		return false, fmt.Errorf("获取所有子节点失败，err:%v", err)
	}
	sort.Sort((sort.StringSlice(childrens)))
	if currentNode == childrens[0] {
		//节点最小，上锁成功
		return true, nil
	}
	//获取上一个节点名
	lastNode := childrens[util.IndexOf(childrens, currentNode)-1]
	return n.WaitLock(lastNode)
}

// 上锁
func (n *ZkLock) Lock() (bool, error) {
	isLock, err := n.TryLock()
	if err != nil {
		return false, err
	}
	if isLock {
		return true, nil
	}
	return false, nil
}

//等待锁
func (n *ZkLock) WaitLock(key string) (bool, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	option := zk.WithEventCallback(func(event zk.Event) {
		fmt.Println("获取到锁了")
		wg.Done()
	})
	conn, _, err := zk.Connect(n.Paths, time.Second*10, option)
	if err != nil {
		return false, fmt.Errorf("等待锁时，开启监听失败,err:%v", err)
	}
	defer conn.Close()
	// 监听上一个节点的状态
	_, _, _, err = conn.ExistsW(fmt.Sprintf("/%s/%s", n.LockPath, key))
	if err != nil {
		return false, fmt.Errorf("等待锁时，获取节点状态失败,err:%v", err)
	}
	wg.Wait()
	return true, nil
}

// 释放锁
func (n *ZkLock) UnLock() error {
	n.c.Close()
	return nil
}

type ZkLockResolver struct {
}

func (n *ZkLockResolver) Resolve(opts ...lock.Option) (lock.LockServer, error) {
	return NewLock(opts...)
}

func init() {
	lock.RegisteLockResolver("zookeeper", &ZkLockResolver{})
}
