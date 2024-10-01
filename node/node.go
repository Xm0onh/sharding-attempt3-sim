package node

import (
	"math/rand"
	"sharding/block"
	"sharding/config"
	"sharding/event"
	"sharding/lottery"
	"sharding/utils"
)

type Node struct {
	ID            int
	IsHonest      bool
	AssignedShard int
	Resources     int
	KnownBlocks   map[int]*block.Block
}

func NewNode(id int) *Node {
	n := &Node{
		ID:            id,
		IsHonest:      true,
		AssignedShard: -1,
		Resources:     1,
		KnownBlocks:   make(map[int]*block.Block),
	}

	if rand.Float64() < config.MaliciousNodeRatio {
		n.IsHonest = false
	}

	return n
}

func (n *Node) ParticipateInLottery(currentTime int64, numShards int) bool {
	win := lottery.WinLottery(n.IsHonest, n.Resources)
	if win {
		n.AssignedShard = lottery.AssignShard(n.ID, currentTime, numShards)
		return true
	}
	return false
}

func (n *Node) CreateBlock(previousBlockID int, currentTime int64) *block.Block {
	blkID := previousBlockID + 1
	blk := block.NewBlock(blkID, n.AssignedShard, n.ID, previousBlockID, currentTime)
	return blk
}

func (n *Node) BroadcastBlock(blk *block.Block, peers []*Node, currentTime int64) []*event.Event {
	events := make([]*event.Event, 0)
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			delay := utils.SimulateNetworkDelay()
			e := &event.Event{
				Timestamp: currentTime + delay,
				Type:      event.MessageEvent,
				NodeID:    peerNode.ID,
				Data:      blk,
			}
			events = append(events, e)
		}
	}
	return events
}

func (n *Node) ProcessMessage(e *event.Event) {
	switch msg := e.Data.(type) {
	case *block.Block:
		n.HandleBlock(msg)
	default:
		// Handle other message types if necessary
	}
}

func (n *Node) HandleBlock(blk *block.Block) {
	if _, exists := n.KnownBlocks[blk.ID]; !exists {
		n.KnownBlocks[blk.ID] = blk
		// The shard's state is managed by the simulation
	}
}
