package goNet

import (
	"github.com/panjf2000/ants"
)

var antsPool *ants.Pool

func initAntsPool() error {
	var err error
	if Opts.PoolSize <= 0 {
		antsPool, err = ants.NewPool(1)
	} else {
		antsPool, err = ants.NewPool(Opts.PoolSize)
	}
	return err
}
