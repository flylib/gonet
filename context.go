package goNet

type Context interface {
	//from session
	Session() Session
	// Message returns the current message to be processed
	Message() interface{}
}

type context struct {
	session Session
	data    interface{}
}

func NewContext(session Session, data interface{}) Context {
	return context{
		session: session,
		data:    data,
	}
}

func (c context) Session() Session {
	return c.session
}

func (c context) Message() interface{} {
	return c.data
}
