package lottery

import (
	"math/rand"
	"sharding/config"
)

func WinLottery(isHonest bool, resources int, currentTime int64, attackStartTime int, attackEndTime int) bool {
	if currentTime >= int64(attackStartTime) && currentTime <= int64(attackEndTime) && !isHonest {
		attempts := resources * config.MaliciousNodeMultiplier
		for i := 0; i < attempts; i++ {
			if rand.Float64() < config.LotteryWinProbability {
				return true
			}
		}
		return false
	}
	return rand.Float64() < config.LotteryWinProbability
}

func AssignShard(nodeID int, timestamp int64, numShards int) int {
	return (nodeID + int(timestamp)) % numShards
}
