package utils

import (
	"math"
	"math/rand"
	cfg "sharding/config"
)

// SimulateNetworkBlockDelay calculates network delay for full block propagation
func SimulateNetworkBlockDelay() float64 {
	// Randomly choose network parameters
	networkDelayMean := cfg.MinNetworkDelayMean + rand.Float64()*(cfg.MaxNetworkDelayMean-cfg.MinNetworkDelayMean)
	networkDelayStd := cfg.MinNetworkDelayStd + rand.Float64()*(cfg.MaxNetworkDelayStd-cfg.MinNetworkDelayStd)
	gossipFanout := cfg.MinGossipFanout + rand.Intn(cfg.MaxGossipFanout-cfg.MinGossipFanout+1)

	// Calculate number of hops in gossip protocol
	numHops := math.Ceil(math.Log(float64(cfg.NumNodes)) / math.Log(float64(gossipFanout)))

	totalDelay := 0.0
	for i := 0.0; i < numHops; i++ {
		// Per-hop latency with jitter
		hopLatency := networkDelayMean +
			rand.NormFloat64()*networkDelayStd/1000.0

		// Transmission delay (size in bits / bandwidth in bps)
		transmissionDelay := (float64(cfg.BlockSize) * 8.0) /
			(float64(cfg.NetworkBandwidth) * 1000000.0)

		totalDelay += hopLatency + transmissionDelay
	}
	return totalDelay
}

// SimulateNetworkBlockHeaderDelay calculates network delay for block header propagation
func SimulateNetworkBlockHeaderDelay() float64 {
	// Randomly choose network parameters
	networkDelayMean := cfg.MinNetworkDelayMean + rand.Float64()*(cfg.MaxNetworkDelayMean-cfg.MinNetworkDelayMean)
	networkDelayStd := cfg.MinNetworkDelayStd + rand.Float64()*(cfg.MaxNetworkDelayStd-cfg.MinNetworkDelayStd)
	gossipFanout := cfg.MinGossipFanout + rand.Intn(cfg.MaxGossipFanout-cfg.MinGossipFanout+1)

	// Calculate number of hops in gossip protocol
	numHops := math.Ceil(math.Log(float64(cfg.NumNodes)) / math.Log(float64(gossipFanout)))

	totalDelay := 0.0
	for i := 0.0; i < numHops; i++ {
		// Per-hop latency with jitter
		hopLatency := networkDelayMean +
			rand.NormFloat64()*networkDelayStd/1000.0

		// Transmission delay (size in bits / bandwidth in bps)
		transmissionDelay := (float64(cfg.BlockHeaderSize) * 8.0) /
			(float64(cfg.NetworkBandwidth) * 1000000.0)

		totalDelay += hopLatency + transmissionDelay
	}

	return totalDelay
}
