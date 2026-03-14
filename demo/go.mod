module demo

go 1.25.0

require (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a
	github.com/flylib/gonet/transport/fastws v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/gonet/transport/gorillaws v0.0.0-20231225121312-42799e3a7e92
	github.com/flylib/goutils/codec/json v0.0.0-20250416114907-55ead4f72e93
	github.com/flylib/pkg/log/builtinlog v0.0.0-20260209033318-13eb902dc7e5
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/fasthttp/websocket v1.5.10 // indirect
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace (
	github.com/flylib/gonet v1.1.4-0.20231101122252-5a97812fe37a => ../
	github.com/flylib/gonet/transport/fastws => ../transport/fastws
	github.com/flylib/gonet/transport/gnet v0.0.0-20231101122252-5a97812fe37a => ../transport/gnet
	github.com/flylib/gonet/transport/gorillaws => ../transport/gorillaws
	github.com/flylib/gonet/transport/quic => ../transport/quic
	github.com/flylib/gonet/transport/udp => ../transport/udp
)
