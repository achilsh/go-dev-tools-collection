module user_demo

go 1.23.0

toolchain go1.24.1

replace user/helloworld => ./stub/user/helloworld

replace student/helloworld => ../student_demo/stub/student/helloworld

require (
	go.uber.org/mock v0.6.0
	student/helloworld v0.0.0-00010101000000-000000000000
	trpc.group/trpc-go/trpc-filter/debuglog v1.0.0
	trpc.group/trpc-go/trpc-filter/recovery v1.0.0
	trpc.group/trpc-go/trpc-go v1.0.3
	trpc.group/trpc-go/trpc-naming-etcd v0.0.0-20250729080045-be5dcfe1cf60
	user/helloworld v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/cenkalti/backoff/v4 v4.1.0 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.1.0 // indirect
	github.com/coreos/pkg v0.0.0-20240122114842-bbd7aa9bf6fb // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/panjf2000/ants/v2 v2.8.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.48.0 // indirect
	go.etcd.io/etcd v3.3.27+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.5.0-alpha.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.0-alpha.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/automaxprocs v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	trpc.group/trpc-go/tnet v1.0.1 // indirect
	trpc.group/trpc/trpc-protocol/pb/go/trpc v1.0.0 // indirect
)
