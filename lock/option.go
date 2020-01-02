package lock

type options struct {
	zkList   []string
	lockPath string
}

type Option func(*options)

func WithZkList(list []string) Option {
	return func(o *options) {
		o.zkList = list
	}
}

func WithLockName(path string) Option {
	return func(o *options) {
		o.lockPath = path
	}
}
