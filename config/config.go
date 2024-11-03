package config

type AttackType int

const (
	TxnSize = 115
)

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
	SimulationTime = 5000 // Total simulation time units
	TimeStep       = 1    // Simulation time step

	// Node parameters
	NumNodes     = 10000
	NumOperators = 20
	// Shard parameters
	NumShards = 1

	// Attack parameters
	MaliciousNodeRatio      = 0.1 // 10% of nodes are malicious
	AttackStartTime         = 20
	AttackEndTime           = 60
	MaliciousNodeMultiplier = 0 // Multiplier for malicious nodes in lottery attempts

	// Lottery parameters
	LotteryWinProbability = 0.001 // Base probability for winning the lottery

	// Block parameters
	BlockProductionInterval = 6                              // Shards produce a block every 6 time steps
	TransactionsPerBlock    = 6500                           // Each block contains 100 transactions
	BlockSize               = TxnSize * TransactionsPerBlock // Block size in bytes
	BlockHeaderSize         = 1000                           // Increased to more realistic size in bytes
	ERHeaderSize            = 1000                           // ER header size in bytes
	ERBodySize              = 33000                          // ER body size in bytes

	// Network simulation parameters
	NetworkBandwidth    = 10    // Network bandwidth in Mbps
	MinNetworkDelayMean = 50.0  // 50ms minimum mean delay
	MaxNetworkDelayMean = 200.0 // 200ms maximum mean delay
	MinNetworkDelayStd  = 10.0  // 10ms minimum standard deviation
	MaxNetworkDelayStd  = 50.0  // 50ms maximum standard deviation
	MinGossipFanout     = 4     // Minimum nodes to gossip to
	MaxGossipFanout     = 8     // Maximum nodes to gossip to
	MaxP2PConnections   = 2
	TimeOut             = 2000 // Timeout for block download in milliseconds

	// Download parameters
	NumBlocksToDownload = 100
)

// InitializeAttackSchedule initializes the attack schedule with both start and end times
func InitializeAttackSchedule() map[int64]AttackType {
	return map[int64]AttackType{
		AttackStartTime: GrindingAttack, // Start Grinding Attack at time step 20
		AttackEndTime:   NoAttack,       // End Grinding Attack at time step 40
	}
}
