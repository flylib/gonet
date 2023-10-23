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
	//net package parser
	netPackageParser INetPackageParser
}

// Default:0 means is no limit
func WithMessageHandler(handler MessageHandler) Option {
	return func(o *option) {
		o.msgHook = handler
	}
}

// Default:0 means is no limit
func WithMaxSessions(max int) Option {
	return func(o *option) {
		o.maxSessionCount = max
	}
}

// Default is runtime.NumCPU(), means no goroutines will be dynamically scaled
func WithPoolMaxRoutines(num int32) Option {
	return func(o *option) {

		o.poolCfg.maxNum = num
	}
}

// allow max idle routines
func WithPoolMaxIdleRoutines(num int32) Option {
	return func(o *option) {
		o.poolCfg.maxIdleNum = num
	}
}

// Default 512,global queue size
func WithGQSize(size int32) Option {
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

// network package paser
func WithNetPackageParser(parser INetPackageParser) Option {
	return func(o *option) {
		o.netPackageParser = parser
	}
}
