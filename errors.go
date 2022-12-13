package gonet

import "errors"

var (
	ErrorSessionClosed = errors.New("session already closed")
	ErrorNotExistMsg   = errors.New("non-existent message")
	ErrorNoTransport   = errors.New("no transmission protocol selected")
)
