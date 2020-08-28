package goNet

import "fmt"

//错误列表
var (
	ErrNotFoundMsg   = fmt.Errorf("The message was not found")
	ErrNotFoundRoute = fmt.Errorf("The Route was not found")
)
