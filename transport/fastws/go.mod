module github.com/flylib/gonet/transport/fastws

go 1.21

require (
	github.com/fasthttp/websocket v1.5.10
	github.com/flylib/gonet v1.1.3
	github.com/valyala/fasthttp v1.58.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.31.0 // indirect
)

replace github.com/flylib/gonet v1.1.3 => ../../../gonet
