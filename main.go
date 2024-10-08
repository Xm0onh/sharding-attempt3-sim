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
		SimulationTime:          config.SimulationTime,
		TimeStep:                config.TimeStep,
		NetworkDelayMean:        config.NetworkDelayMean,
		NetworkDelayStd:         config.NetworkDelayStd,
		MaliciousNodeRatio:      config.MaliciousNodeRatio,
		LotteryWinProbability:   config.LotteryWinProbability,
		MaliciousNodeMultiplier: config.MaliciousNodeMultiplier,
		BlockProductionInterval: config.BlockProductionInterval,
		TransactionsPerBlock:    config.TransactionsPerBlock,
		AttackSchedule:          config.InitializeAttackSchedule(),
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
