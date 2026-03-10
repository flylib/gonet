package gonet

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
)

// Option configures a Context via NewAppContext.
type Option func(*config)

func WithEventHandler(handler IEventHandler) Option {
	return func(c *config) { c.eventHandler = handler }
}

// WithMaxSessions limits the maximum number of concurrent sessions (0 = unlimited).
func WithMaxSessions(max int) Option {
	return func(c *config) { c.maxSessionCount = max }
}

// WithPoolMaxRoutines sets the hard cap on worker goroutines (0 = unlimited).
func WithPoolMaxRoutines(num int32) Option {
	return func(c *config) { c.poolCfg.maxNum = num }
}

// WithPoolMaxIdleRoutines sets the initial/idle worker count (default: NumCPU).
func WithPoolMaxIdleRoutines(num int32) Option {
	return func(c *config) { c.poolCfg.maxIdleNum = num }
}

// WithGQSize sets the global message queue buffer size (default: 64).
func WithGQSize(size int32) Option {
	return func(c *config) { c.poolCfg.queueSize = size }
}

// WithNetPackager sets a custom network packet encoder/decoder.
func WithNetPackager(p INetPackager) Option {
	return func(c *config) { c.INetPackager = p }
}

// MustWithCodec sets the message codec (required).
func MustWithCodec(c codec.ICodec) Option {
	return func(cfg *config) { cfg.ICodec = c }
}

// MustWithLogger sets the logger (required).
func MustWithLogger(l ilog.ILogger) Option {
	return func(cfg *config) { cfg.ILogger = l }
}
