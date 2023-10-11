module github.com/flylib/gonet/transport/ws

go 1.21.1

require (
	github.com/flylib/gonet v1.1.3
	github.com/gorilla/websocket v1.5.0
)

require (
	github.com/flylib/goutils/logger/log v0.0.0-20231010072000-719eaf1057f5 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/flylib/gonet v1.1.3 => ../../../gonet
