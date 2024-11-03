package metrics

import (
	"fmt"
	"io"
	"os"
	"sharding/config"
	"sharding/node"
	"sharding/shard"
)

type NetworkMetrics struct {
	BlockBroadcastDelays map[int][]float64
	BlockHeaderDelays    []float64
	BlockDownloadDelays  map[int][]float64
	AverageBlockDelay    map[int]float64
	AverageHeaderDelay   float64
	AverageDownloadDelay map[int]float64
}

type ShardMetrics struct {
	HonestNodes     int
	MaliciousNodes  int
	HonestBlocks    int
	MaliciousBlocks int
}

type TimeWindowMetrics struct {
	TotalEvents             int64
	AverageResponseTime     float64
	ErrorRate               float64
	TotalBlocks             int
	TotalTransactions       int
	MaliciousShardRotations int
	NetworkMetrics          NetworkMetrics
	ShardStats              map[int]*ShardMetrics
}

type MetricsCollector struct {
	CurrentMetrics TimeWindowMetrics
	Logs           []string
}

type SimulationResponse struct {
	TransactionSize      int                  `json:"transaction_size_bytes"`
	TransactionsPerBlock int                  `json:"transactions_per_block"`
	BlockSize            int                  `json:"block_size_kb"`
	BlockProduction      map[int]ShardStats   `json:"block_production"`
	NetworkMetrics       NetworkStatsResponse `json:"network_metrics"`
	Performance          PerformanceStats     `json:"performance"`
}

type ShardStats struct {
	MaliciousBlocks int `json:"malicious_blocks"`
	HonestBlocks    int `json:"honest_blocks"`
	TotalBlocks     int `json:"total_blocks"`
}

type NetworkStatsResponse struct {
	BlockBroadcastDelays map[int]float64 `json:"block_broadcast_delays_ms"`
	BlockHeaderDelay     float64         `json:"block_header_delay_ms"`
	BlockDownloadDelays  map[int]float64 `json:"block_download_delays_ms"`
}

type PerformanceStats struct {
	TPS float64 `json:"transactions_per_second"`
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		CurrentMetrics: TimeWindowMetrics{
			ShardStats: make(map[int]*ShardMetrics),
			NetworkMetrics: NetworkMetrics{
				BlockBroadcastDelays: make(map[int][]float64),
				BlockDownloadDelays:  make(map[int][]float64),
				AverageBlockDelay:    make(map[int]float64),
				AverageDownloadDelay: make(map[int]float64),
			},
		},
		Logs: make([]string, 0),
	}
}

func (mc *MetricsCollector) Collect(
	timestamp int64,
	shards map[int]*shard.Shard,
	nodes map[int]*node.Node,
	blockDelays map[int][]int64,
	headerDelays map[int][]int64,
	downloadDelays map[int][]int64,
	logs []string,
	maliciousRotations int,
) {
	// Process network delays
	for shardID, delays := range blockDelays {
		if mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays[shardID] == nil {
			mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays[shardID] = make([]float64, 0)
		}
		for _, delay := range delays {
			mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays[shardID] = append(
				mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays[shardID],
				float64(delay),
			)
		}
	}

	for _, delays := range headerDelays {
		for _, delay := range delays {
			mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays = append(
				mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays,
				float64(delay),
			)
		}
	}

	// Update download delays processing
	for shardID, delays := range downloadDelays {
		if mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays[shardID] == nil {
			mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays[shardID] = make([]float64, 0)
		}
		for _, delay := range delays {
			mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays[shardID] = append(
				mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays[shardID],
				float64(delay),
			)
		}
	}

	// Reset shard statistics for this collection
	mc.CurrentMetrics.ShardStats = make(map[int]*ShardMetrics)

	// Reset total blocks counter
	mc.CurrentMetrics.TotalBlocks = 0

	// Update shard statistics
	for shardID, s := range shards {
		stats := &ShardMetrics{}
		mc.CurrentMetrics.ShardStats[shardID] = stats

		// Count honest and malicious nodes
		honestNodes := s.GetHonestNodes()
		maliciousNodes := s.GetMaliciousNodes()
		stats.HonestNodes = len(honestNodes)
		stats.MaliciousNodes = len(maliciousNodes)

		// Count blocks in the shard
		stats.HonestBlocks = 0
		stats.MaliciousBlocks = 0
		for _, block := range s.Blocks {
			if block.IsMalicious {
				stats.MaliciousBlocks++
			} else {
				stats.HonestBlocks++
			}
		}

		// Update total blocks count
		mc.CurrentMetrics.TotalBlocks += stats.HonestBlocks + stats.MaliciousBlocks
	}

	mc.CurrentMetrics.MaliciousShardRotations += maliciousRotations
	mc.Logs = append(mc.Logs, logs...)
}

func (mc *MetricsCollector) calculateAverages() {
	// Calculate broadcast delays per shard
	for shardID, delays := range mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays {
		if len(delays) > 0 {
			sum := 0.0
			for _, d := range delays {
				sum += d
			}
			mc.CurrentMetrics.NetworkMetrics.AverageBlockDelay[shardID] = sum / float64(len(delays))
		}
	}

	if len(mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays) > 0 {
		sum := 0.0
		for _, d := range mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays {
			sum += d
		}
		mc.CurrentMetrics.NetworkMetrics.AverageHeaderDelay = sum / float64(len(mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays))
	}

	// Calculate download delays per shard
	for shardID, delays := range mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays {
		if len(delays) > 0 {
			sum := 0.0
			for _, d := range delays {
				sum += d
			}
			mc.CurrentMetrics.NetworkMetrics.AverageDownloadDelay[shardID] = sum / float64(len(delays))
		}
	}
}

func (mc *MetricsCollector) GenerateReport() error {
	// Calculate averages before generating report
	mc.calculateAverages()

	f, err := os.Create("simulation_report.txt")
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "=== Simulation Report ===")
	writeTimeWindowMetrics(f, "Simulation Metrics", mc.CurrentMetrics)

	// Write logs
	fmt.Fprintln(f, "=== Event Logs ===")
	for _, log := range mc.Logs {
		fmt.Fprintln(f, log)
	}

	return nil
}

// Helper function to calculate percentage
func CalculatePercentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func writeTimeWindowMetrics(w io.Writer, title string, metrics TimeWindowMetrics) {
	fmt.Fprintf(w, "   Size of each Transaction in bytes: ~%d\n", 100)
	fmt.Fprintf(w, "   Number of transactions per block: %d\n", config.TransactionsPerBlock)
	fmt.Fprintf(w, "   Size of each Block in kilo bytes: %d\n", config.BlockSize/1000)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s:\n", title)

	// Add block production statistics
	fmt.Fprintf(w, "\nBlock Production Statistics:\n")
	for shardID, stats := range metrics.ShardStats {
		totalBlocks := stats.HonestBlocks + stats.MaliciousBlocks
		fmt.Fprintf(w, "=== Shard %d ===\n", shardID)
		fmt.Fprintf(w, "  Shard %d: %d malicious blocks\n", shardID, stats.MaliciousBlocks)
		fmt.Fprintf(w, "  Shard %d: %d honest blocks\n", shardID, stats.HonestBlocks)
		fmt.Fprintf(w, "  Shard %d: %d blocks\n", shardID, totalBlocks)

	}
	fmt.Fprintf(w, "\n") // Add spacing

	// Original metrics
	// fmt.Fprintf(w, "Simulation Metrics:\n")
	// fmt.Fprintf(w, "  Total Events: %d\n", metrics.TotalEvents)
	// fmt.Fprintf(w, "  Average Response Time: %.2fms\n", metrics.AverageResponseTime)
	// fmt.Fprintf(w, "  Error Rate: %.2f%%\n", metrics.ErrorRate*100)

	// Network metrics
	fmt.Fprintf(w, "\nNetwork Metrics:\n")
	fmt.Fprintf(w, "  Average Block Broadcast Delay per Shard:\n")
	for shardID, avgDelay := range metrics.NetworkMetrics.AverageBlockDelay {
		fmt.Fprintf(w, "    Shard %d: %.2fms\n", shardID, avgDelay)
	}
	fmt.Fprintf(w, "  Average Block Header Delay: %.2fms\n", metrics.NetworkMetrics.AverageHeaderDelay)
	fmt.Fprintf(w, "  Average Block Download Delay per Shard:\n")
	for shardID, avgDelay := range metrics.NetworkMetrics.AverageDownloadDelay {
		fmt.Fprintf(w, "    Shard %d: %.2fms\n", shardID, avgDelay)
	}
	fmt.Fprintf(w, "\n")

	// Add TPS calculation
	totalBlocks := 0
	for _, stats := range metrics.ShardStats {
		totalBlocks += stats.HonestBlocks
	}
	totalTransactions := totalBlocks * config.TransactionsPerBlock
	fmt.Println("Total txn:", totalTransactions)
	tps := float64(totalTransactions) / float64(config.SimulationTime)
	fmt.Fprintf(w, "Performance Metrics:\n")
	fmt.Fprintf(w, "  Transactions Per Second (TPS): %.2f\n\n", tps)
}

func (mc *MetricsCollector) GetSimulationResponse() SimulationResponse {
	mc.calculateAverages()

	response := SimulationResponse{
		TransactionSize:      100, // Fixed value from the report
		TransactionsPerBlock: config.TransactionsPerBlock,
		BlockSize:            config.BlockSize / 1000,
		BlockProduction:      make(map[int]ShardStats),
		NetworkMetrics: NetworkStatsResponse{
			BlockBroadcastDelays: mc.CurrentMetrics.NetworkMetrics.AverageBlockDelay,
			BlockHeaderDelay:     mc.CurrentMetrics.NetworkMetrics.AverageHeaderDelay,
			BlockDownloadDelays:  mc.CurrentMetrics.NetworkMetrics.AverageDownloadDelay,
		},
	}

	// Calculate total blocks and populate shard stats
	totalBlocks := 0
	for shardID, stats := range mc.CurrentMetrics.ShardStats {
		totalBlocks += stats.HonestBlocks
		response.BlockProduction[shardID] = ShardStats{
			MaliciousBlocks: stats.MaliciousBlocks,
			HonestBlocks:    stats.HonestBlocks,
			TotalBlocks:     stats.HonestBlocks + stats.MaliciousBlocks,
		}
	}

	// Calculate TPS
	totalTransactions := totalBlocks * config.TransactionsPerBlock
	tps := float64(totalTransactions) / float64(config.SimulationTime)
	response.Performance = PerformanceStats{
		TPS: tps,
	}

	return response
}
