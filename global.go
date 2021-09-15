package gonet

import (
	"reflect"
	"sync"
)

var (
	mgr           manager
	transportType TransportType //传输协议类型
	globalLock    sync.RWMutex
)

func init() {
	mgr = manager{
		pool: sync.Pool{
			New: func() interface{} {
				return reflect.New(transportType).Interface()
			},
		},
	}
}
