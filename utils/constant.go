package utils

const (
	MaliciousNodeRatio      = 0.1   // 10% of nodes are malicious
	LotteryWinProbability   = 0.001 // Adjusted for large number of nodes
	NetworkDelayMean        = 5     // Average network delay in time units
	NetworkDelayStd         = 2     // Standard deviation of network delay
	MaliciousNodeMultiplier = 2     // Multiplier for malicious nodes in lottery attempts
)
