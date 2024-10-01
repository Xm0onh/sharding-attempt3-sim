// attack/attack.go

package attack

import (
	"fmt"
	"sharding/config"
	"sharding/event"
	"sharding/node"
	"sharding/shard"
)

// ExecuteAttack performs the specified attack based on the AttackType.
// It appends attack-related logs to the simulation's AttackLogs.
func ExecuteAttack(atkType config.AttackType, currentTime int64, nodes map[int]*node.Node, shards map[int]*shard.Shard, eq *event.EventQueue, cfg config.Config, attackLogs *[]string) {
	switch atkType {
	case config.GrindingAttack:
		performGrindingAttack(currentTime, nodes, eq, cfg, attackLogs)
	case config.NoAttack:
		// stopGrindingAttack(currentTime, nodes, shards, eq, cfg, attackLogs)
	default:
		// Unknown attack type
		log := fmt.Sprintf("[Attack] Unknown attack type: %v at time %d", atkType, currentTime)
		*attackLogs = append(*attackLogs, log)
	}
}

// performGrindingAttack schedules additional LotteryEvents for malicious nodes to increase their shard assignments.
func performGrindingAttack(currentTime int64, nodes map[int]*node.Node, eq *event.EventQueue, cfg config.Config, attackLogs *[]string) {
	log := fmt.Sprintf("[Attack] Performing Grinding Attack at time %d", currentTime)
	fmt.Println("Current Time:", currentTime)
	*attackLogs = append(*attackLogs, log)

	// for _, n := range nodes {
	// 	if !n.IsHonest {
	// 		// Schedule additional LotteryEvents based on the multiplier
	// 		for i := 0; i < cfg.MaliciousNodeMultiplier; i++ {
	// 			extraLottery := &event.Event{
	// 				Timestamp: currentTime, // Immediate participation
	// 				Type:      event.LotteryEvent,
	// 				NodeID:    n.ID,
	// 			}
	// 			heap.Push(eq, extraLottery)
	// 		}
	// 		log := fmt.Sprintf("[Attack] Malicious Node %d scheduled %d additional LotteryEvents", n.ID, cfg.MaliciousNodeMultiplier)
	// 		*attackLogs = append(*attackLogs, log)
	// 	}
	// }
}

// // stopGrindingAttack logs the termination of the Grinding Attack.
// // Currently, no specific action is required to stop the attack.
// // This function serves as a placeholder for potential future implementations.
func stopGrindingAttack(currentTime int64, nodes map[int]*node.Node, shards map[int]*shard.Shard, eq *event.EventQueue, cfg config.Config, attackLogs *[]string) {
	log := fmt.Sprintf("[Attack] Stopping Grinding Attack at time %d", currentTime)
	*attackLogs = append(*attackLogs, log)
}
