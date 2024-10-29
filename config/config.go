package config

type AttackType int

const (
	NoAttack AttackType = iota
	GrindingAttack
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
	BlockSize               int
	BlockHeaderSize         int
	ERHeaderSize            int
	ERBodySize              int
	NetworkBandwidth        int64
	GossipFanout            int
}

const (
	// Simulation parameters
	SimulationTime = 120 // Total simulation time units
	TimeStep       = 1   // Simulation time step

	// Node parameters
	NumNodes = 100

	// Shard parameters
	NumShards = 1

	// Attack parameters
	MaliciousNodeRatio      = 0.1 // 10% of nodes are malicious
	AttackStartTime         = 20
	AttackEndTime           = 60
	MaliciousNodeMultiplier = 1000 // Multiplier for malicious nodes in lottery attempts

	// Lottery parameters
	LotteryWinProbability = 0.01 // Base probability for winning the lottery

	// Block parameters
	BlockProductionInterval = 6     // Shards produce a block every 6 time steps
	TransactionsPerBlock    = 100   // Each block contains 100 transactions
	BlockSize               = 10000 // Block size in bytes
	BlockHeaderSize         = 1000  // Increased to more realistic size in bytes
	ERHeaderSize            = 1000  // ER header size in bytes
	ERBodySize              = 33000 // ER body size in bytes

	// Network parameters
	NetworkDelayMean = 100 // Updated to milliseconds for more realism
	NetworkDelayStd  = 50  // Updated standard deviation
	NetworkBandwidth = 10  // Network bandwidth in Mbps
	GossipFanout     = 4   // Each node forwards to 4 peers by default
)

// InitializeAttackSchedule initializes the attack schedule with both start and end times
func InitializeAttackSchedule() map[int64]AttackType {
	return map[int64]AttackType{
		AttackStartTime: GrindingAttack, // Start Grinding Attack at time step 20
		AttackEndTime:   NoAttack,       // End Grinding Attack at time step 40
	}
}
