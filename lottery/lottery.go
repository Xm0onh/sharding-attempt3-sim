package lottery

import (
	"math/rand"
	"sharding/config"
)

func WinLottery(isHonest bool, resources int) bool {
	if isHonest {
		return rand.Float64() < config.LotteryWinProbability
	} else {
		attempts := resources * config.MaliciousNodeMultiplier
		for i := 0; i < attempts; i++ {
			if rand.Float64() < config.LotteryWinProbability {
				return true
			}
		}
		return false
	}
}

func AssignShard(nodeID int, timestamp int64, numShards int) int {
	return (nodeID + int(timestamp)) % numShards
}
