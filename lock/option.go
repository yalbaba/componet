package lock

type Options struct {
	ZkList   []string
	LockPath string
}

type Option func(*Options)

func WithZkList(list []string) Option {
	return func(o *Options) {
		o.ZkList = list
	}
}

func WithLockName(path string) Option {
	return func(o *Options) {
		o.LockPath = path
	}
}
