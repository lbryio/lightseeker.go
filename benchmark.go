package main

import (
	"time"
)

type benchmark interface {
	Start()
	Rate() float64
	Summary() string
	Timing() time.Duration
}
