package lock

import "fmt"

const EPHEMERAL_SEQUENTIAL = 3

var lockResolvers map[string]LockResolver

type LockServer interface {
	TryLock() (bool, error)
	Lock() (bool, error)
	WaitLock(last string) (bool, error)
	UnLock()
}

//根据名称获取锁对象
func GetLockServer(proto string, opts ...Option) (LockServer, error) {
	resolver, ok := lockResolvers[proto]
	if !ok {
		fmt.Errorf("没有适配器")
	}
	return resolver.Resolve(opts...), nil
}

type LockResolver interface {
	Resolve(opts ...Option) LockServer
}

func RegisteLockResolver(proto string, resolver LockResolver) error {
	if resolver == nil {
		fmt.Errorf("注入的适配器不可为空")
	}
	if _, ok := lockResolvers[proto]; ok {
		fmt.Errorf("该适配器已经注入")
	}
	lockResolvers[proto] = resolver
	return nil
}
