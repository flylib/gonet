package gonet

type IEventHandler interface {
	OnConnect(ISession)
	OnClose(ISession, error)
	OnMessage(message Message)
	OnError(ISession, error)
}
