package shard

import (
	"sharding/block"
)

type Shard struct {
	ID     int
	Blocks []*block.Block
}

func NewShard(id int) *Shard {
	s := &Shard{
		ID:     id,
		Blocks: make([]*block.Block, 0),
	}
	return s
}

func (s *Shard) AddBlock(blk *block.Block) {
	s.Blocks = append(s.Blocks, blk)
}

func (s *Shard) LatestBlockID() int {
	if len(s.Blocks) == 0 {
		return 0 // Genesis block ID
	}
	return s.Blocks[len(s.Blocks)-1].ID
}
