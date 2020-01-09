module zklock

go 1.13

replace github.com/samuel/go-zookeeper/zk => ../github.com/samuel/go-zookeeper/zk

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
)
