// main.go

package main

import (
	"fmt"
	"sharding/config"
	"sharding/simulation"
)

func main() {
	// Initialize simulation parameters
	cfg := config.Config{
		NumNodes:                config.NumNodes,
		NumShards:               config.NumShards,
		NumOperators:            config.NumOperators,
		SimulationTime:          config.SimulationTime,
		TimeStep:                config.TimeStep,
		MaliciousNodeRatio:      config.MaliciousNodeRatio,
		LotteryWinProbability:   config.LotteryWinProbability,
		MaliciousNodeMultiplier: config.MaliciousNodeMultiplier,
		BlockProductionInterval: config.BlockProductionInterval,
		TransactionsPerBlock:    config.TransactionsPerBlock,
		AttackSchedule:          config.InitializeAttackSchedule(),
		BlockSize:               config.BlockSize,
		BlockHeaderSize:         config.BlockHeaderSize,
		ERHeaderSize:            config.ERHeaderSize,
		ERBodySize:              config.ERBodySize,
		NetworkBandwidth:        config.NetworkBandwidth,
		MinNetworkDelayMean:     config.MinNetworkDelayMean,
		MaxNetworkDelayMean:     config.MaxNetworkDelayMean,
		MinNetworkDelayStd:      config.MinNetworkDelayStd,
		MaxNetworkDelayStd:      config.MaxNetworkDelayStd,
		MinGossipFanout:         config.MinGossipFanout,
		MaxGossipFanout:         config.MaxGossipFanout,
		MaxP2PConnections:       config.MaxP2PConnections,
		TimeOut:                 config.TimeOut,
		NumBlocksToDownload:     config.NumBlocksToDownload,
	}

	// Create a new simulation instance
	sim := simulation.NewSimulation(cfg)

	// Run the simulation
	fmt.Println("Simulation started.")
	sim.Run()
	fmt.Println("Simulation completed.")

	// Generate metrics report
	err := sim.Metrics.GenerateReport()
	if err != nil {
		fmt.Printf("Error generating metrics report: %v\n", err)
	} else {
		fmt.Println("Metrics report generated.")
	}
}
