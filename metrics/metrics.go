package metrics

import (
	"fmt"
	"os"
	"sharding/config"
	"sharding/node"
	"sharding/shard"
)

type ShardMetrics struct {
	HonestNodes     int
	MaliciousNodes  int
	HonestBlocks    int
	MaliciousBlocks int
}

type MetricsData struct {
	Timestamp           int64
	TotalBlocks         int
	TotalTransactions   int
	Throughput          float64
	AverageNetworkDelay float64
	Latency             float64
	ShardStats          map[int]*ShardMetrics
}

type MetricsCollector struct {
	Data []MetricsData
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		Data: make([]MetricsData, 0),
	}
}

func (mc *MetricsCollector) Collect(timestamp int64, shards map[int]*shard.Shard, nodes map[int]*node.Node, networkDelays []int64) {
	md := MetricsData{
		Timestamp:           timestamp,
		TotalBlocks:         0,
		TotalTransactions:   0,
		AverageNetworkDelay: 0,
		Latency:             0,
		ShardStats:          make(map[int]*ShardMetrics),
	}

	totalDelay := int64(0)
	totalEvents := len(networkDelays)

	// Initialize ShardStats
	for shardID := range shards {
		md.ShardStats[shardID] = &ShardMetrics{}
	}

	// Count nodes per shard
	for _, n := range nodes {
		if n.IsAssignedToShard() {
			shardID := n.AssignedShard
			if n.IsHonest {
				md.ShardStats[shardID].HonestNodes++
			} else {
				md.ShardStats[shardID].MaliciousNodes++
			}
		}
	}

	// Count blocks per shard
	for shardID, s := range shards {
		shardMetrics := md.ShardStats[shardID]
		md.TotalBlocks += len(s.Blocks)
		md.TotalTransactions += len(s.Blocks) * config.TransactionsPerBlock

		for _, blk := range s.Blocks {
			if blk.IsMalicious {
				shardMetrics.MaliciousBlocks++
			} else {
				shardMetrics.HonestBlocks++
			}
		}
	}

	for _, delay := range networkDelays {
		totalDelay += delay
	}

	if totalEvents > 0 {
		md.AverageNetworkDelay = float64(totalDelay) / float64(totalEvents)
	}

	// Calculate Throughput (TPS)
	timeElapsed := timestamp
	if timeElapsed > 0 {
		md.Throughput = float64(md.TotalTransactions) / float64(timeElapsed)
	}

	// Calculate Latency
	md.Latency = float64(config.BlockProductionInterval) + md.AverageNetworkDelay

	mc.Data = append(mc.Data, md)
}

func (mc *MetricsCollector) GenerateReport() error {
	file, err := os.Create("metrics_report.txt")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("Metrics Report:\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for _, md := range mc.Data {
		_, err := file.WriteString(fmt.Sprintf("Time: %d, Total Blocks: %d, Total Transactions: %d, TPS: %.2f, Avg Latency: %.2f units\n",
			md.Timestamp, md.TotalBlocks, md.TotalTransactions, md.Throughput, md.Latency))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		for shardID, stats := range md.ShardStats {
			_, err := file.WriteString(fmt.Sprintf("  Shard %d: Honest Nodes: %d, Malicious Nodes: %d, Honest Blocks: %d, Malicious Blocks: %d\n",
				shardID, stats.HonestNodes, stats.MaliciousNodes, stats.HonestBlocks, stats.MaliciousBlocks))
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}

	return nil
}
