package lock

type Options struct {
	Paths    []string
	LockPath string
	Address  string
	Password string
	DB       int
	TimeOut  int
}

type Option func(o *Options)

func WithZkList(list []string) Option {
	return func(o *Options) {
		o.Paths = list
	}
}

func WithLockName(path string) Option {
	return func(o *Options) {
		o.LockPath = path
	}
}
func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func WithPassword(pwd string) Option {
	return func(o *Options) {
		o.Password = pwd
	}
}

func WithTimeOut(time int) Option {
	return func(o *Options) {
		o.TimeOut = time
	}
}
