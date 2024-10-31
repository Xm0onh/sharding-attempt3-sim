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
	BlockBroadcastDelays []float64
	BlockHeaderDelays    []float64
	BlockDownloadDelays  []float64
	AverageBlockDelay    float64
	AverageHeaderDelay   float64
	AverageDownloadDelay float64
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

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		CurrentMetrics: TimeWindowMetrics{
			ShardStats: make(map[int]*ShardMetrics),
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
	for _, delays := range blockDelays {
		for _, delay := range delays {
			mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays = append(
				mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays,
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

	for _, delays := range downloadDelays {
		for _, delay := range delays {
			mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays = append(
				mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays,
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
	if len(mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays) > 0 {
		sum := 0.0
		for _, d := range mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays {
			sum += d
		}
		mc.CurrentMetrics.NetworkMetrics.AverageBlockDelay = sum / float64(len(mc.CurrentMetrics.NetworkMetrics.BlockBroadcastDelays))
	}

	if len(mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays) > 0 {
		sum := 0.0
		for _, d := range mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays {
			sum += d
		}
		mc.CurrentMetrics.NetworkMetrics.AverageHeaderDelay = sum / float64(len(mc.CurrentMetrics.NetworkMetrics.BlockHeaderDelays))
	}

	if len(mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays) > 0 {
		sum := 0.0
		for _, d := range mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays {
			sum += d
		}
		mc.CurrentMetrics.NetworkMetrics.AverageDownloadDelay = sum / float64(len(mc.CurrentMetrics.NetworkMetrics.BlockDownloadDelays))
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
	// fmt.Fprintln(f, "=== Event Logs ===")
	// for _, log := range mc.Logs {
	// 	fmt.Fprintln(f, log)
	// }

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
	fmt.Fprintf(w, "  Average Block Broadcast Delay: %.2fms\n", metrics.NetworkMetrics.AverageBlockDelay)
	fmt.Fprintf(w, "  Average Block Header Delay: %.2fms\n", metrics.NetworkMetrics.AverageHeaderDelay)
	fmt.Fprintf(w, "  Average Block Download Delay: %.2fms\n\n", metrics.NetworkMetrics.AverageDownloadDelay)

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
