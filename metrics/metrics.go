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
	Timestamp                       int64
	BlocksThisStep                  int
	TransactionsThisStep            int
	MaliciousShardRotationsThisStep int // New field
	Throughput                      float64
	AverageNetworkDelay             float64
	Latency                         float64
	ShardStats                      map[int]*ShardMetrics
	Logs                            []string
}

type MetricsCollector struct {
	Data []MetricsData
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		Data: make([]MetricsData, 0),
	}
}

// Collect gathers metrics at each time step, including attack logs and malicious shard rotations.
func (mc *MetricsCollector) Collect(timestamp int64, shards map[int]*shard.Shard, nodes map[int]*node.Node, networkDelays []int64, Logs []string, maliciousShardRotations int) {
	md := MetricsData{
		Timestamp:                       timestamp,
		BlocksThisStep:                  0,
		TransactionsThisStep:            0,
		MaliciousShardRotationsThisStep: maliciousShardRotations,
		AverageNetworkDelay:             0,
		Latency:                         0,
		ShardStats:                      make(map[int]*ShardMetrics),
		Logs:                            make([]string, len(Logs)),
	}

	// Copy attack logs to avoid mutation
	copy(md.Logs, Logs)

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
				if !blk.IsMalicious {
					md.TransactionsThisStep += config.TransactionsPerBlock
					md.ShardStats[shardID].HonestBlocks++
				} else {
					md.ShardStats[shardID].MaliciousBlocks++
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

	// Variables to accumulate malicious shard rotations
	var preAttackRotations, duringAttackRotations, postAttackRotations int

	// Define attack window
	attackStart := config.AttackStartTime
	attackEnd := config.AttackEndTime

	// Accumulate total blocks produced in each shard
	totalBlocksPerShard := make(map[int]int)
	totalHonestBlocksPerShard := make(map[int]int)
	totalMaliciousBlocksPerShard := make(map[int]int)
	for _, md := range mc.Data {
		for shardID, stats := range md.ShardStats {
			totalBlocksPerShard[shardID] += stats.HonestBlocks + stats.MaliciousBlocks
			totalHonestBlocksPerShard[shardID] += stats.HonestBlocks
			totalMaliciousBlocksPerShard[shardID] += stats.MaliciousBlocks
		}
	}

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
		for _, log := range md.Logs {
			_, err := file.WriteString(fmt.Sprintf("  %s\n", log))
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}

		// Write malicious shard rotations, if any
		if md.MaliciousShardRotationsThisStep > 0 {
			_, err := file.WriteString(fmt.Sprintf("  Malicious Shard Rotations This Step: %d\n", md.MaliciousShardRotationsThisStep))
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}

		// Categorize TPS and Malicious Rotations based on attack window
		if md.Timestamp < int64(attackStart) {
			preAttackTPS += md.Throughput
			preAttackCount++
			preAttackRotations += md.MaliciousShardRotationsThisStep
		} else if md.Timestamp >= int64(attackStart) && md.Timestamp <= int64(attackEnd) {
			duringAttackTPS += md.Throughput
			duringAttackCount++
			duringAttackRotations += md.MaliciousShardRotationsThisStep
		} else {
			postAttackTPS += md.Throughput
			postAttackCount++
			postAttackRotations += md.MaliciousShardRotationsThisStep
		}
		// Add a separator for readability
		_, err = file.WriteString("\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	// Write the total number of blocks produced in each shard
	for shardID, totalBlocks := range totalBlocksPerShard {
		_, err := file.WriteString(fmt.Sprintf("Total Blocks Produced in Shard %d: %d\n", shardID, totalBlocks))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	// Write the total number of honest and malicious blocks in each shard
	for shardID := range totalBlocksPerShard {
		_, err := file.WriteString(fmt.Sprintf("Total Honest Blocks in Shard %d: %d\n", shardID, totalHonestBlocksPerShard[shardID]))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		_, err = file.WriteString(fmt.Sprintf("Total Malicious Blocks in Shard %d: %d\n", shardID, totalMaliciousBlocksPerShard[shardID]))
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

	// Calculate total rotations for each period
	totalPreAttackRotations := preAttackRotations
	totalDuringAttackRotations := duringAttackRotations
	totalPostAttackRotations := postAttackRotations

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

	// Append rotations summary
	_, err = file.WriteString(fmt.Sprintf("\nTotal Successful Shard Rotations by Malicious Nodes before Attack (Time < %d): %d\n", attackStart, totalPreAttackRotations))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Total Successful Shard Rotations by Malicious Nodes during Attack (%d <= Time <= %d): %d\n", attackStart, attackEnd, totalDuringAttackRotations))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Total Successful Shard Rotations by Malicious Nodes after Attack (Time > %d): %d\n", attackEnd, totalPostAttackRotations))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// Analyze the effect of Grinding Attack on TPS
	_, err = file.WriteString("\nEffect of Grinding Attack on TPS:\n")
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

	// Analyze the effect of Grinding Attack on shard rotations
	_, err = file.WriteString("\nEffect of Grinding Attack on Shard Rotations:\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if totalDuringAttackRotations > totalPreAttackRotations {
		_, err = file.WriteString("The Grinding Attack led to more successful shard rotations by malicious nodes during the attack period.\n")
	} else if totalDuringAttackRotations < totalPreAttackRotations {
		_, err = file.WriteString("The Grinding Attack did not increase the number of successful shard rotations by malicious nodes during the attack period.\n")
	} else {
		_, err = file.WriteString("The Grinding Attack had no significant effect on the number of successful shard rotations by malicious nodes during the attack period.\n")
	}

	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
