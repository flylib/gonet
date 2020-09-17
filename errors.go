package goNet

import "fmt"

//错误列表
var (
	ErrNotFoundMsg   = fmt.Errorf("The message was not found")
	ErrNotFoundActor = fmt.Errorf("The Actor was not found")
)
