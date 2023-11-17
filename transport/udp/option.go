package udp

import "time"

type Option func(*option)

type option struct {
	//specifies the duration for the handshake to complete.Default is 5 second
	HandshakeTimeout time.Duration
	mtu              int
}

// set maximum transmission unit
func WithMaximumTransmissionUnit(unit int) Option {
	return func(o *option) {
		o.mtu = unit
	}
}
