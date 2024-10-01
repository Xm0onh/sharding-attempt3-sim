package main

import (
	"fmt"
	"sharding/config"
	"sharding/simulation"
)

func main() {
	// Initialize simulation parameters
	cfg := config.Config{
		NumNodes:         config.NumNodes,
		NumShards:        config.NumShards,
		SimulationTime:   config.SimulationTime,
		TimeStep:         config.TimeStep,
		NetworkDelayMean: config.NetworkDelayMean,
		NetworkDelayStd:  config.NetworkDelayStd,
		AttackStartTime:  config.AttackStartTime,
		AttackEndTime:    config.AttackEndTime,
		AttackType:       config.NoAttack, // Change as needed
	}

	// Create a new simulation instance
	sim := simulation.NewSimulation(cfg)

	// Run the simulation
	fmt.Println("Simulation started.")
	sim.Run()
	fmt.Println("Simulation completed.")

	// Generate metrics report
	sim.Metrics.GenerateReport()
}
