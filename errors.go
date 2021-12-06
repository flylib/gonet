package gonet

import "errors"

var (
	ErrorSessionClosed = errors.New("Session already closed")
	ErrorNotExistMsg   = errors.New("Non-existent message")
	ErrorNoTransport   = errors.New("No transmission protocol selected")
)
