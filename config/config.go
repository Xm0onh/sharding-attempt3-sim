// config/config.go

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
	MaliciousNodeRatio      float64
	LotteryWinProbability   float64
	MaliciousNodeMultiplier int
	AttackSchedule          map[int64]AttackType
}

const (
	NumNodes                = 10000
	NumShards               = 10
	SimulationTime          = 120   // Total simulation time units
	TimeStep                = 1     // Simulation time step
	NetworkDelayMean        = 5     // Average network delay in time units
	NetworkDelayStd         = 2     // Standard deviation of network delay
	MaliciousNodeRatio      = 0.1   // 10% of nodes are malicious
	LotteryWinProbability   = 0.001 // Base probability for winning the lottery
	AttackStartTime         = 20
	AttackEndTime           = 60
	MaliciousNodeMultiplier = 1000 // Multiplier for malicious nodes in lottery attempts
	BlockProductionInterval = 6    // Shards produce a block every 6 time steps
	TransactionsPerBlock    = 100  // Each block contains 100 transactions
)

// InitializeAttackSchedule initializes the attack schedule with both start and end times
func InitializeAttackSchedule() map[int64]AttackType {
	return map[int64]AttackType{
		AttackStartTime: GrindingAttack, // Start Grinding Attack at time step 20
		AttackEndTime:   NoAttack,       // End Grinding Attack at time step 40
	}
}
