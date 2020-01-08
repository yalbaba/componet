package redis

import (
	"time"
	"zklock/lock"

	"github.com/robfig/cron"
	"github.com/go-redis/redis"
)

type RedisLock struct {
	value   string
	ch      chan string
	cron 
	client  *redis.Client
	OptConf *lock.Options
}

func NewLock(opts ...lock.Option) (*RedisLock, error) {
	// 使用uuid
	lock := &RedisLock{value: "value"}
	for _, o := range opts {
		o(lock.OptConf)
	}

	lock.client = redis.NewClient(&redis.Options{
		Addr:     lock.OptConf.Address,
		Password: lock.OptConf.Password,
		DB:       lock.OptConf.DB,
	})

	_, err := lock.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return lock, nil
}

func (l *RedisLock) TryLock() (bool, error) {
	// 获取锁（锁值应为请求的唯一标识，标识这是某次请求的值）
	value, err := l.client.SetNX(l.OptConf.LockPath, l.value, time.Second*time.Duration(l.OptConf.TimeOut)).Result()
	if err != nil {
		return false, err
	}
	if value {
		// 开启定时任务每隔10秒监控这个key

		return true, nil
	}

	return false, nil
}

func (l *RedisLock) Lock() (bool, error) {
	return false, nil
}

func (l *RedisLock) WaitLock(last string) (bool, error) {
	return false, nil
}

func (l *RedisLock) UnLock() error {
	val, err := l.client.Get(l.OptConf.LockPath).Result()
	if err != nil {
		return err
	}
	//判断本次请求value是否等于当前锁的value
	if val != l.value {
		l.client.Del(l.OptConf.LockPath)
		l.client.Close()
	}
	return nil
}

type RedisLockResolver struct{}

func (n *RedisLockResolver) Resolve(opts ...lock.Option) (lock.ZkLockServer, error) {
	return NewLock(opts...)
}

func init() {
	lock.RegisteLockResolver("redis", &RedisLockResolver{})
}

// c := cron.New()
//     spec := "*/5 * * * * ?"
//     c.AddFunc(spec, func() {
//         i++
//         log.Println("cron running:", i)
//     })
//     c.Start()
