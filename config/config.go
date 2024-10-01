package config

type AttackType int

const (
	NoAttack AttackType = iota
	GrindingAttack
	// Add other attack types as needed
)

type Config struct {
	NumNodes                int
	NumShards               int
	SimulationTime          int64
	TimeStep                int64
	NetworkDelayMean        int64
	NetworkDelayStd         int64
	AttackStartTime         int64
	AttackEndTime           int64
	AttackType              AttackType
	BlockProductionInterval int64
	TransactionsPerBlock    int
}

const (
	NumNodes                = 100
	NumShards               = 10
	SimulationTime          = 100  // Total simulation time units
	TimeStep                = 1    // Simulation time step
	NetworkDelayMean        = 5    // Average network delay in time units
	NetworkDelayStd         = 2    // Standard deviation of network delay
	MaliciousNodeRatio      = 0.1  // 10% of nodes are malicious
	LotteryWinProbability   = 0.01 // Increased probability for demonstration
	AttackStartTime         = 20
	AttackEndTime           = 40
	MaliciousNodeMultiplier = 2   // Multiplier for malicious nodes in lottery attempts
	BlockProductionInterval = 5   // Shards produce a block every 5 time steps
	TransactionsPerBlock    = 100 // Each block contains 100 transactions
)
