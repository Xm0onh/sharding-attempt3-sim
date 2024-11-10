// shard/shard.go

package shard

import (
	"fmt"
	"sharding/block"
	"sharding/node"
)

type Shard struct {
	ID     int
	Blocks []*block.Block
	Nodes  map[int]*node.Node
}

func NewShard(id int) *Shard {
	s := &Shard{
		ID:     id,
		Blocks: make([]*block.Block, 0),
		Nodes:  make(map[int]*node.Node),
	}
	return s
}

func (s *Shard) AddBlock(blk *block.Block) {
	for _, b := range s.Blocks {
		if b.ID == blk.ID {
			return
		}
	}
	s.Blocks = append(s.Blocks, blk)
}

func (s *Shard) LatestBlockID() int {
	if len(s.Blocks) == 0 {
		return 0 // Genesis block ID
	}
	return s.Blocks[len(s.Blocks)-1].ID
}

func (s *Shard) AddNode(n *node.Node) {
	s.Nodes[n.ID] = n
	// fmt.Printf("[Shard %d] Node %d added. Total Nodes: %d\n", s.ID, n.ID, len(s.Nodes))
}

func (s *Shard) RemoveNode(nodeID int) {
	if _, exists := s.Nodes[nodeID]; exists {
		delete(s.Nodes, nodeID)
		// fmt.Printf("[Shard %d] Node %d removed. Total Nodes: %d\n", s.ID, nodeID, len(s.Nodes))
	}
}

func (s *Shard) GetNodes() []*node.Node {
	nodes := make([]*node.Node, 0, len(s.Nodes))
	for _, n := range s.Nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (s *Shard) GetHonestNodes() []*node.Node {
	honest := []*node.Node{}
	for _, n := range s.Nodes {
		if n.IsHonest {
			honest = append(honest, n)
		}
	}
	return honest
}

func (s *Shard) GetMaliciousNodes() []*node.Node {
	malicious := []*node.Node{}
	for _, n := range s.Nodes {
		if !n.IsHonest {
			malicious = append(malicious, n)
		}
	}
	return malicious
}

func (s *Shard) IsolateNode(nodeID int) {
	if _, exists := s.Nodes[nodeID]; exists {
		s.RemoveNode(nodeID)
		fmt.Printf("[Shard %d] Node %d has been isolated.\n", s.ID, nodeID)
	}
}

func (s *Shard) GetBlock(id int) *block.Block {
	return s.Blocks[id]
}

func (s *Shard) GetLatestBlockID() int {
	latestID := 0
	for _, block := range s.Blocks {
		if block.ID > latestID {
			latestID = block.ID
		}
	}
	return latestID
}
