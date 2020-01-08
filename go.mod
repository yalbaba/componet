module zklock

go 1.13

replace github.com/samuel/go-zookeeper/zk => ../github.com/samuel/go-zookeeper/zk

replace github.com/gofrs/uuid => ../github.com/gofrs/uuid

require (
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gofrs/uuid v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/samuel/go-zookeeper/zk v0.0.0-00010101000000-000000000000
)
