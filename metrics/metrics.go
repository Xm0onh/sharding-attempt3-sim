package metrics

import (
	"fmt"
	"os"
	"sharding/node"
	"sharding/shard"
)

type MetricsData struct {
	Timestamp           int64
	TotalBlocks         int
	MaliciousBlocks     int
	AverageNetworkDelay float64
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
		MaliciousBlocks:     0,
		AverageNetworkDelay: 0,
	}

	totalDelay := int64(0)
	for _, s := range shards {
		blockCount := len(s.Blocks)
		md.TotalBlocks += blockCount

		for _, blk := range s.Blocks {
			if !nodes[blk.ProducerID].IsHonest {
				md.MaliciousBlocks++
			}
		}
	}

	for _, delay := range networkDelays {
		totalDelay += delay
	}

	if len(networkDelays) > 0 {
		md.AverageNetworkDelay = float64(totalDelay) / float64(len(networkDelays))
	}

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
		_, err := file.WriteString(fmt.Sprintf("Time: %d, Total Blocks: %d, Malicious Blocks: %d, Avg Network Delay: %.2f units\n",
			md.Timestamp, md.TotalBlocks, md.MaliciousBlocks, md.AverageNetworkDelay))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
