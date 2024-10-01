package block

type Block struct {
	ID           int
	ShardID      int
	ProducerID   int
	PreviousHash int
	Timestamp    int64
}

func NewBlock(id, shardID, producerID, previousHash int, timestamp int64) *Block {
	return &Block{
		ID:           id,
		ShardID:      shardID,
		ProducerID:   producerID,
		PreviousHash: previousHash,
		Timestamp:    timestamp,
	}
}
