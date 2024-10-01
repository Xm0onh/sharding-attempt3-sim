package utils

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SimulateNetworkDelay() int64 {
	delay := rand.NormFloat64()*float64(NetworkDelayStd) + float64(NetworkDelayMean)
	if delay < 1 {
		delay = 1
	}
	return int64(math.Round(delay))
}
