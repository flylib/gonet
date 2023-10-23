package gonet

type Option func(o *option)

type option struct {
	//Message callback processing
	msgHook         MessageHandler
	maxSessionCount int
	//routine pool config
	poolCfg poolConfig
	//message codec
	codec ICodec
	log   ILogger
}

// Default:0 means is no limit
func WithMaxSessions(max int) Option {
	return func(o *option) {
		o.maxSessionCount = max
	}
}

// Default is runtime.NumCPU(), means no goroutines will be dynamically scaled
func WithPoolMaxRoutines(num uint32) Option {
	return func(o *option) {

		o.poolCfg.maxNum = num
	}
}

// allow max idle routines
func WithPoolMaxIdleRoutines(num uint32) Option {
	return func(o *option) {
		o.poolCfg.maxIdleNum = num
	}
}

// Default 512,global queue size
func WithGQSize(size uint32) Option {
	return func(o *option) {
		o.poolCfg.queueSize = size
	}
}

// Default json codec, message codec
func WithMessageCodec(codec ICodec) Option {
	return func(o *option) {
		o.codec = codec
	}
}

// set logger
func WithLogger(l ILogger) Option {
	return func(o *option) {
		o.log = l
	}
}
