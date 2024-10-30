package config

type AttackType int

const (
	NoAttack AttackType = iota
	GrindingAttack
)

type Config struct {
	NumNodes                int
	NumOperators            int
	NumShards               int
	SimulationTime          int64
	TimeStep                int64
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
	MinNetworkDelayMean     float64
	MaxNetworkDelayMean     float64
	MinNetworkDelayStd      float64
	MaxNetworkDelayStd      float64
	MinGossipFanout         int
	MaxGossipFanout         int
	MaxP2PConnections       int
	TimeOut                 int64
	NumBlocksToDownload     int
}

const (
	// Simulation parameters
	SimulationTime = 120 // Total simulation time units
	TimeStep       = 1   // Simulation time step

	// Node parameters
	NumNodes     = 100_000
	NumOperators = 10
	// Shard parameters
	NumShards = 1

	// Attack parameters
	MaliciousNodeRatio      = 0.1 // 10% of nodes are malicious
	AttackStartTime         = 20
	AttackEndTime           = 60
	MaliciousNodeMultiplier = 10 // Multiplier for malicious nodes in lottery attempts

	// Lottery parameters
	LotteryWinProbability = 0.001 // Base probability for winning the lottery

	// Block parameters
	BlockProductionInterval = 6     // Shards produce a block every 6 time steps
	TransactionsPerBlock    = 100   // Each block contains 100 transactions
	BlockSize               = 10000 // Block size in bytes
	BlockHeaderSize         = 1000  // Increased to more realistic size in bytes
	ERHeaderSize            = 1000  // ER header size in bytes
	ERBodySize              = 33000 // ER body size in bytes

	// Network simulation parameters
	NetworkBandwidth    = 10    // Network bandwidth in Mbps
	MinNetworkDelayMean = 50.0  // 50ms minimum mean delay
	MaxNetworkDelayMean = 200.0 // 200ms maximum mean delay
	MinNetworkDelayStd  = 10.0  // 10ms minimum standard deviation
	MaxNetworkDelayStd  = 50.0  // 50ms maximum standard deviation
	MinGossipFanout     = 4     // Minimum nodes to gossip to
	MaxGossipFanout     = 8     // Maximum nodes to gossip to
	MaxP2PConnections   = 1
	TimeOut             = 2000 // Timeout for block download in milliseconds

	// Download parameters
	NumBlocksToDownload = 10
)

// InitializeAttackSchedule initializes the attack schedule with both start and end times
func InitializeAttackSchedule() map[int64]AttackType {
	return map[int64]AttackType{
		AttackStartTime: GrindingAttack, // Start Grinding Attack at time step 20
		AttackEndTime:   NoAttack,       // End Grinding Attack at time step 40
	}
}
