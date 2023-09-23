package gonet

// /////////////////////////////
// ///    Option Func   ////////
// ////////////////////////////
type Option func(*Context) error

func MaxSessions(max int) Option {
	return func(o *Context) error {
		o.maxSessionCount = max
		return nil
	}
}

func WorkerPoolMaxSize(max int) Option {
	return func(o *Context) error {
		o.maxWorkerPoolSize = max
		return nil
	}
}

// cache for messages
func WithMessageCache(cache IMessageCache) Option {
	return func(o *Context) error {
		o.msgCache = cache
		return nil
	}
}

// message codec,default is json codec
func WithMessageCodec(codec ICodec) Option {
	return func(o *Context) error {
		o.codec = codec
		return nil
	}
}
