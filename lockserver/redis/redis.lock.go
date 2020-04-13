package redis

import (
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"time"
	"zklock/lockserver"
)

type RedisLock struct {
	request string
	ch      chan int
	client  *redis.Client
	OptConf *lockserver.Options
}

func NewLock(opts ...lockserver.Option) (*RedisLock, error) {
	// 使用uuid
	uuID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	lock := &RedisLock{
		request: uuID.String(),
		ch:      make(chan int),
	}
	for _, o := range opts {
		o(lock.OptConf)
	}

	lock.client = redis.NewClient(&redis.Options{
		Addr:     lock.OptConf.Address,
		Password: lock.OptConf.Password,
		DB:       lock.OptConf.DB,
	})

	_, err = lock.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return lock, nil
}

// 执行上锁
func (l *RedisLock) TryLock() (bool, error) {
	// 获取锁（锁值应为请求的唯一标识，标识这是某次请求的值）
	value, err := l.client.SetNX(l.OptConf.LockPath, l.request, time.Second*time.Duration(l.OptConf.TimeOut)).Result()
	if err != nil {
		return false, err
	}
	if value {
		return true, nil
	}

	return false, nil
}

// 开始上锁
func (l *RedisLock) Lock() (bool, error) {
	for {
		isSuccess, err := l.TryLock()
		if err != nil {
			return false, err
		}
		if isSuccess {
			l.ch <- 1
		}
	}

	select {
	case <-l.ch:
		return true, nil
	}
}

func (l *RedisLock) WaitLock(key string) (bool, error) {

	return false, nil
}

// 释放锁
func (l *RedisLock) UnLock() error {
	val, err := l.client.Get(l.OptConf.LockPath).Result()
	if err != nil {
		return err
	}
	//判断本次请求value是否等于当前锁的request
	if val != l.request {
		// 此处也存在安全问题，例如程序刚执行到这里过期了，接下来删除的还是其他进程的锁
		l.client.Del(l.OptConf.LockPath)
		l.client.Close()
	}
	return nil
}

type RedisLockResolver struct{}

func (n *RedisLockResolver) Resolve(opts ...lockserver.Option) (lockserver.LockServer, error) {
	return NewLock(opts...)
}

func init() {
	lockserver.RegisteLockResolver("redis", &RedisLockResolver{})
}
