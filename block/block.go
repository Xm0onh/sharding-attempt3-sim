package block

type Block struct {
	ID           int
	ShardID      int
	ProducerID   int
	PreviousHash int
	Timestamp    int64
	IsMalicious  bool
}

type BlockHeader struct {
	ID         int
	ShardID    int
	ProducerID int
	PreviousID int
	Timestamp  int64
}

func NewBlockHeader(id, shardID, producerID, previousID int, timestamp int64) *BlockHeader {
	return &BlockHeader{
		ID:         id,
		ShardID:    shardID,
		ProducerID: producerID,
		PreviousID: previousID,
		Timestamp:  timestamp,
	}
}

func NewBlock(id, shardID, producerID, previousHash int, timestamp int64) *Block {
	return &Block{
		ID:           id,
		ShardID:      shardID,
		ProducerID:   producerID,
		PreviousHash: previousHash,
		Timestamp:    timestamp,
		IsMalicious:  false, // Default to false
	}
}
