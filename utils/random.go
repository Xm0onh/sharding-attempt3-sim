package utils

import (
	"math/rand"
	cfg "sharding/config"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SimulateNetworkBlockDelay() float64 {
	delay := rand.NormFloat64()*float64(NetworkDelayStd) + float64(NetworkDelayMean)
	return delay
}

func SimulateNetworkBlockHeaderDelay() float64 {
	// Base delay using normal distribution
	baseDelay := rand.NormFloat64()*float64(NetworkDelayStd) + float64(NetworkDelayMean)

	// Scale delay based on block header size (assuming 1KB = 1 time unit scaling factor)
	sizeScalingFactor := float64(cfg.BlockHeaderSize) / 1000.0

	return baseDelay * sizeScalingFactor
}
