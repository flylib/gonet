module demo

go 1.25.0

require (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/gnet v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/gonet/transport/quic v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/gonet/transport/udp v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/gonet/transport/ws v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/goutils/codec/json v0.0.0-20250416114907-55ead4f72e93
	github.com/flylib/pkg/log/builtinlog v0.0.0-20260209033318-13eb902dc7e5
)

require (
	github.com/flylib/goutils/sync/spinlock v0.0.0-20250416114907-55ead4f72e93 // indirect
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20260302011040-a15ffb7f9dcc // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.28.1 // indirect
	github.com/panjf2000/ants/v2 v2.11.5 // indirect
	github.com/panjf2000/gnet/v2 v2.9.7 // indirect
	github.com/quic-go/qtls-go1-20 v0.4.1 // indirect
	github.com/quic-go/quic-go v0.59.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a => ../
	github.com/flylib/gonet/transport/gnet v0.0.0-20231101122252-5a97812fe37a => ../transport/gnet
	github.com/flylib/gonet/transport/quic => ../transport/quic
	github.com/flylib/gonet/transport/udp => ../transport/udp
	github.com/flylib/gonet/transport/ws => ../transport/ws
)
