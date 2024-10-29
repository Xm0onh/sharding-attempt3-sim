package utils

import (
	"math"
	"math/rand"
	cfg "sharding/config"
)

// SimulateNetworkBlockDelay calculates network delay for full block propagation
func SimulateNetworkBlockDelay() float64 {
	// Calculate number of hops in gossip protocol
	numHops := math.Ceil(math.Log(float64(cfg.NumNodes)) / math.Log(float64(cfg.GossipFanout)))

	totalDelay := 0.0
	for i := 0.0; i < numHops; i++ {
		// Per-hop latency with jitter
		hopLatency := cfg.NetworkLatencyBase +
			rand.NormFloat64()*float64(cfg.NetworkDelayStd)/1000.0

		// Transmission delay (size in bits / bandwidth in bps)
		transmissionDelay := (float64(cfg.BlockSize) * 8.0) /
			(float64(cfg.NetworkBandwidth) * 1000000.0)

		totalDelay += hopLatency + transmissionDelay
	}
	return totalDelay * 1000.0 // Convert to milliseconds
}

// SimulateNetworkBlockHeaderDelay calculates network delay for block header propagation
func SimulateNetworkBlockHeaderDelay() float64 {
	// Calculate number of hops in gossip protocol
	numHops := math.Ceil(math.Log(float64(cfg.NumNodes)) / math.Log(float64(cfg.GossipFanout)))

	totalDelay := 0.0
	for i := 0.0; i < numHops; i++ {
		// Per-hop latency with jitter
		hopLatency := cfg.NetworkLatencyBase +
			rand.NormFloat64()*float64(cfg.NetworkDelayStd)/1000.0

		// Transmission delay (size in bits / bandwidth in bps)
		transmissionDelay := (float64(cfg.BlockHeaderSize) * 8.0) /
			(float64(cfg.NetworkBandwidth) * 1000000.0)

		totalDelay += hopLatency + transmissionDelay
	}

	return totalDelay * 1000.0 // Convert to milliseconds
}
