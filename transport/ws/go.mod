module github.com/flylib/gonet/transport/ws

go 1.18


require (
	github.com/flylib/gonet v1.1.3
	github.com/gorilla/websocket v1.5.0
)

require (
	github.com/flylib/goutils/codec/json v0.0.0-20231012070911-2cf6c2bcb71d // indirect
	github.com/flylib/goutils/logger v0.0.0-20231023014531-4f50a5871c60 // indirect
	github.com/flylib/goutils/logger/log v0.0.0-20231023014531-4f50a5871c60 // indirect
	github.com/flylib/goutils/sync/spinlock v0.0.0-20231019075452-c7b2623472a2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/flylib/gonet v1.1.3 => ../../../gonet
