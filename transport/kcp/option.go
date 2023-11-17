package kcp

import (
	"time"
)

type Option func(*option)

type option struct {
	//specifies the duration for the handshake to complete.Default is 5 second
	HandshakeTimeout time.Duration
	//PBKDF的全称是Password-Based Key Derivation Function，简单的说，PBKDF就是一个密码衍生的工具。
	//PBKDF2和PBKDF1主要是用来防止密码暴力破解的，所以在设计中加入了对算力的自动调整，从而抵御暴力破解的可能性。
	PBKDF2Password, PBKDF2Salt string
}

func WithHandshakeTimeout(t time.Duration) Option {
	return func(option *option) {
		option.HandshakeTimeout = t
	}
}

func WithPBKDF2(pwd, salt string) Option {
	return func(o *option) {
		o.PBKDF2Password = pwd
		o.PBKDF2Salt = salt
	}
}
