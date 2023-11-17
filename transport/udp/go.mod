module github.com/flylib/gonet/transport/udp

go 1.18

require github.com/flylib/gonet v1.1.3

require (
	github.com/flylib/goutils/container v0.0.0-20231115102727-0f7df9653a51 // indirect
	github.com/flylib/goutils/sync/spinlock v0.0.0-20231019075452-c7b2623472a2 // indirect
	github.com/flylib/interface v0.0.0-20231101042444-4c3b4b8d0e0d // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
)

replace github.com/flylib/gonet v1.1.3 => ../../../gonet
