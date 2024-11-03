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
	http.HandleFunc("/simulate-with-config", handleSimulationWithConfig)

	// Start HTTP server
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

// UserConfig matches the frontend configuration structure
type UserConfig struct {
	NumNodes                int     `json:"numNodes"`
	NumShards               int     `json:"numShards"`
	NumOperators            int     `json:"numOperators"`
	SimulationTime          int64   `json:"simulationTime"`
	TimeStep                int64   `json:"timeStep"`
	MaliciousNodeRatio      float64 `json:"maliciousNodeRatio"`
	LotteryWinProbability   float64 `json:"lotteryWinProbability"`
	MaliciousNodeMultiplier int     `json:"maliciousNodeMultiplier"`
	BlockProductionInterval int64   `json:"blockProductionInterval"`
	TransactionsPerBlock    int     `json:"transactionsPerBlock"`
	BlockSize               int     `json:"blockSize"`
	BlockHeaderSize         int     `json:"blockHeaderSize"`
	ERHeaderSize            int     `json:"erHeaderSize"`
	ERBodySize              int     `json:"erBodySize"`
	NetworkBandwidth        int64   `json:"networkBandwidth"`
	MinNetworkDelayMean     float64 `json:"minNetworkDelayMean"`
	MaxNetworkDelayMean     float64 `json:"maxNetworkDelayMean"`
	MinNetworkDelayStd      float64 `json:"minNetworkDelayStd"`
	MaxNetworkDelayStd      float64 `json:"maxNetworkDelayStd"`
	MinGossipFanout         int     `json:"minGossipFanout"`
	MaxGossipFanout         int     `json:"maxGossipFanout"`
	MaxP2PConnections       int     `json:"maxP2PConnections"`
	TimeOut                 int64   `json:"timeOut"`
	NumBlocksToDownload     int     `json:"numBlocksToDownload"`
	AttackStartTime         int64   `json:"attackStartTime"`
	AttackEndTime           int64   `json:"attackEndTime"`
}

func handleSimulationWithConfig(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulationMutex.Lock()
	defer simulationMutex.Unlock()

	// Parse the user configuration
	var userConfig UserConfig
	if err := json.NewDecoder(r.Body).Decode(&userConfig); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse configuration: %v", err), http.StatusBadRequest)
		return
	}
	// Printing the received config from the user
	fmt.Println("Received config from the user:", userConfig)
	fmt.Println("Num of blocks to download:", userConfig.NumBlocksToDownload)
	// Initialize metrics collector
	metricsCollector = metrics.NewMetricsCollector()

	// Convert user config to simulation config
	cfg := config.Config{
		NumNodes:                userConfig.NumNodes,
		NumShards:               userConfig.NumShards,
		NumOperators:            userConfig.NumOperators,
		SimulationTime:          userConfig.SimulationTime,
		TimeStep:                userConfig.TimeStep,
		MaliciousNodeRatio:      userConfig.MaliciousNodeRatio,
		LotteryWinProbability:   userConfig.LotteryWinProbability,
		MaliciousNodeMultiplier: userConfig.MaliciousNodeMultiplier,
		BlockProductionInterval: userConfig.BlockProductionInterval,
		TransactionsPerBlock:    userConfig.TransactionsPerBlock,
		BlockSize:               userConfig.BlockSize,
		BlockHeaderSize:         userConfig.BlockHeaderSize,
		ERHeaderSize:            userConfig.ERHeaderSize,
		ERBodySize:              userConfig.ERBodySize,
		NetworkBandwidth:        userConfig.NetworkBandwidth,
		MinNetworkDelayMean:     userConfig.MinNetworkDelayMean,
		MaxNetworkDelayMean:     userConfig.MaxNetworkDelayMean,
		MinNetworkDelayStd:      userConfig.MinNetworkDelayStd,
		MaxNetworkDelayStd:      userConfig.MaxNetworkDelayStd,
		MinGossipFanout:         userConfig.MinGossipFanout,
		MaxGossipFanout:         userConfig.MaxGossipFanout,
		MaxP2PConnections:       userConfig.MaxP2PConnections,
		TimeOut:                 userConfig.TimeOut,
		NumBlocksToDownload:     userConfig.NumBlocksToDownload,
		AttackSchedule: map[int64]config.AttackType{
			userConfig.AttackStartTime: config.GrindingAttack,
			userConfig.AttackEndTime:   config.NoAttack,
		},
	}

	// Create a new simulation instance with metrics collector
	sim := simulation.NewSimulation(cfg, metricsCollector)

	// Run the simulation
	fmt.Println("Simulation started with custom configuration.")
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
