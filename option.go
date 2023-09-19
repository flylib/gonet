package gonet

// /////////////////////////////
// ///    Option Func   ////////
// ////////////////////////////
type HandlerFunc func(*Context) error

func MaxSessions(max int) HandlerFunc {
	return func(o *Context) error {
		o.maxSessionCount = max
		return nil
	}
}

func WorkerPoolMaxSize(max int) HandlerFunc {
	return func(o *Context) error {
		o.maxWorkerPoolSize = max
		return nil
	}
}

// cache for messages
func WithMessageCache(cache IEventCache) HandlerFunc {
	return func(o *Context) error {
		o.msgCache = cache
		return nil
	}
}
