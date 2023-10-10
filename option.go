package gonet

type Option func(*AppContext) error

func MaxSessions(max int) Option {
	return func(o *AppContext) error {
		o.maxSessionCount = max
		return nil
	}
}

func WorkerPoolMaxSize(max int) Option {
	return func(o *AppContext) error {
		o.maxWorkerPoolSize = max
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
