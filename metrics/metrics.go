// metrics.go

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
	Timestamp            int64
	BlocksThisStep       int
	TransactionsThisStep int
	Throughput           float64
	AverageNetworkDelay  float64
	Latency              float64
	ShardStats           map[int]*ShardMetrics
	AttackLogs           []string
}

type MetricsCollector struct {
	Data []MetricsData
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		Data: make([]MetricsData, 0),
	}
}

// Collect gathers metrics at each time step, including attack logs.
func (mc *MetricsCollector) Collect(timestamp int64, shards map[int]*shard.Shard, nodes map[int]*node.Node, networkDelays []int64, attackLogs []string) {
	md := MetricsData{
		Timestamp:            timestamp,
		BlocksThisStep:       0,
		TransactionsThisStep: 0,
		AverageNetworkDelay:  0,
		Latency:              0,
		ShardStats:           make(map[int]*ShardMetrics),
		AttackLogs:           make([]string, len(attackLogs)),
	}

	// Copy attack logs to avoid mutation
	copy(md.AttackLogs, attackLogs)

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
			if shardID < 0 {
				continue
			}
			if _, exists := md.ShardStats[shardID]; !exists {
				md.ShardStats[shardID] = &ShardMetrics{}
			}
			if n.IsHonest {
				md.ShardStats[shardID].HonestNodes++
			} else {
				md.ShardStats[shardID].MaliciousNodes++
			}
		}
	}

	// Count blocks and transactions per shard for this time step
	for shardID, s := range shards {
		for _, blk := range s.Blocks {
			if blk.Timestamp == timestamp {
				md.BlocksThisStep++
				md.TransactionsThisStep += config.TransactionsPerBlock
				if blk.IsMalicious {
					md.ShardStats[shardID].MaliciousBlocks++
				} else {
					md.ShardStats[shardID].HonestBlocks++
				}
			}
		}
	}

	// Calculate average network delay
	for _, delay := range networkDelays {
		totalDelay += delay
	}

	if totalEvents > 0 {
		md.AverageNetworkDelay = float64(totalDelay) / float64(totalEvents)
	}

	// Calculate Throughput (TPS) for this step
	md.Throughput = float64(md.TransactionsThisStep) / float64(config.TimeStep)

	// Calculate Latency
	md.Latency = float64(config.BlockProductionInterval) + md.AverageNetworkDelay

	mc.Data = append(mc.Data, md)
}

// GenerateReport writes the collected metrics to a report file, including attack logs and a summary analysis.
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

	// Variables to accumulate TPS for analysis
	var preAttackTPS, duringAttackTPS, postAttackTPS float64
	var preAttackCount, duringAttackCount, postAttackCount int

	// Define attack window
	attackStart := config.AttackStartTime
	attackEnd := config.AttackEndTime

	for _, md := range mc.Data {
		// Write global metrics
		_, err := file.WriteString(fmt.Sprintf("Time: %d, Blocks This Step: %d, Transactions This Step: %d, TPS: %.2f, Avg Latency: %.2f units\n",
			md.Timestamp, md.BlocksThisStep, md.TransactionsThisStep, md.Throughput, md.Latency))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		// Write shard-specific metrics
		for shardID, stats := range md.ShardStats {
			_, err := file.WriteString(fmt.Sprintf("  Shard %d: Honest Nodes: %d, Malicious Nodes: %d, Honest Blocks: %d, Malicious Blocks: %d\n",
				shardID, stats.HonestNodes, stats.MaliciousNodes, stats.HonestBlocks, stats.MaliciousBlocks))
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}

		// Write attack logs, if any
		for _, log := range md.AttackLogs {
			_, err := file.WriteString(fmt.Sprintf("  %s\n", log))
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}

		// Categorize TPS based on attack window
		if md.Timestamp < int64(attackStart) {
			preAttackTPS += md.Throughput
			preAttackCount++
		} else if md.Timestamp >= int64(attackStart) && md.Timestamp <= int64(attackEnd) {
			duringAttackTPS += md.Throughput
			duringAttackCount++
		} else {
			postAttackTPS += md.Throughput
			postAttackCount++
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	// Calculate average TPS for each period
	avgPreAttackTPS := 0.0
	if preAttackCount > 0 {
		avgPreAttackTPS = preAttackTPS / float64(preAttackCount)
	}

	avgDuringAttackTPS := 0.0
	if duringAttackCount > 0 {
		avgDuringAttackTPS = duringAttackTPS / float64(duringAttackCount)
	}

	avgPostAttackTPS := 0.0
	if postAttackCount > 0 {
		avgPostAttackTPS = postAttackTPS / float64(postAttackCount)
	}

	// Append summary analysis
	_, err = file.WriteString("Summary Analysis:\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Average TPS before Grinding Attack (Time < %d): %.2f\n", attackStart, avgPreAttackTPS))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Average TPS during Grinding Attack (%d <= Time <= %d): %.2f\n", attackStart, attackEnd, avgDuringAttackTPS))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Average TPS after Grinding Attack (Time > %d): %.2f\n", attackEnd, avgPostAttackTPS))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// Analyze the effect of Grinding Attack
	_, err = file.WriteString("\nEffect of Grinding Attack:\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if avgDuringAttackTPS > avgPreAttackTPS {
		_, err = file.WriteString("The Grinding Attack increased the TPS during the attack period.\n")
	} else if avgDuringAttackTPS < avgPreAttackTPS {
		_, err = file.WriteString("The Grinding Attack decreased the TPS during the attack period.\n")
	} else {
		_, err = file.WriteString("The Grinding Attack had no significant effect on the TPS during the attack period.\n")
	}

	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
