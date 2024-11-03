// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sharding/config"
	"sharding/metrics"
	"sharding/simulation"
	"sync"
)

var (
	metricsCollector *metrics.MetricsCollector
	simulationMutex  sync.Mutex
)

func main() {
	// Setup HTTP routes
	http.HandleFunc("/simulate", handleSimulation)

	// Start HTTP server
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleSimulation(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulationMutex.Lock()
	defer simulationMutex.Unlock()

	// Initialize metrics collector
	metricsCollector = metrics.NewMetricsCollector()

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

	// Create a new simulation instance with metrics collector
	sim := simulation.NewSimulation(cfg, metricsCollector)

	// Run the simulation
	fmt.Println("Simulation started.")
	sim.Run()
	fmt.Println("Simulation completed.")

	// Generate metrics report
	err := metricsCollector.GenerateReport()
	if err != nil {
		fmt.Printf("Error generating metrics report: %v\n", err)
	} else {
		fmt.Println("Metrics report generated.")
	}

	// Get the metrics response
	response := metricsCollector.GetSimulationResponse()

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
