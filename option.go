package gonet

import "errors"

type Option func(*AppContext) error

func MaxSessions(max int) Option {
	return func(o *AppContext) error {
		o.maxSessionCount = max
		return nil
	}
}

// 0 means no goroutines will be dynamically scaled
func MaxWorkers(num uint32) Option {
	return func(o *AppContext) error {
		if num > 1024 {
			return errors.New("Setting too many workers is not allowed")
		}
		o.workerOptions = append(o.workerOptions, maxWorkingGoroutines(int32(num)))
		return nil
	}
}

func MaxIdleWorkers(num uint32) Option {
	return func(o *AppContext) error {
		if num == 0 {
			return errors.New("zero workers is not allowed")
		}
		o.workerOptions = append(o.workerOptions, maxIdleGoroutines(int32(num)))
		return nil
	}
}

func GlobalMessageQueueSize(num uint32) Option {
	return func(o *AppContext) error {
		o.workerOptions = append(o.workerOptions, setQueueSize(int(num)))
		return nil
	}
}

// message codec,default is json codec
func WithMessageCodec(codec ICodec) Option {
	return func(o *AppContext) error {
		o.codec = codec
		return nil
	}
}

// set logger
func Logger(l ILogger) Option {
	return func(o *AppContext) error {
		o.ILogger = l
		return nil
	}
}
