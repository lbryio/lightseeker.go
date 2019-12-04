package main

import (
	"time"
)

type benchmark interface {
	Start()
	InstantholdRate() float64
	ThresholdRate() float64
	WholesomeRate() float64
	Summary() string
	Timing() time.Duration
}
