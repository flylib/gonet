module demo

go 1.18

require (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/gnet v0.0.0-20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/quic v0.0.0-20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/udp v0.0.0-20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/ws v0.0.0-20231101122252-5a97812fe37a
	github.com/flylib/goutils/codec/json v0.0.0-20231026110424-19dfbb98ff56
	github.com/flylib/pkg/log/builtinlog v0.0.0-20231031025337-eee45d016863
)

require (
	github.com/flylib/goutils/sync/spinlock v0.0.0-20231019075452-c7b2623472a2 // indirect
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.9.5 // indirect
	github.com/panjf2000/gnet/v2 v2.3.3 // indirect
	github.com/quic-go/qtls-go1-20 v0.4.1 // indirect
	github.com/quic-go/quic-go v0.40.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/mock v0.3.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a => ../
	github.com/flylib/gonet/transport/gnet v0.0.0-20231101122252-5a97812fe37a => ../transport/gnet
	github.com/flylib/gonet/transport/quic => ../transport/quic
	github.com/flylib/gonet/transport/udp => ../transport/udp
	github.com/flylib/gonet/transport/ws => ../transport/ws
)
