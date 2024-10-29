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
		AssignedShard: -1, // Unassigned initially
		Resources:     1,
		KnownBlocks:   make(map[int]*block.Block),
	}

	if rand.Float64() < config.MaliciousNodeRatio {
		n.IsHonest = false
	}

	return n
}

func (n *Node) ParticipateInLottery(currentTime int64, numShards int) (bool, int) {
	if n.IsAssignedToShard() {
		return false, -1 // Already assigned to a shard
	}

	win := lottery.WinLottery(n.IsHonest, 1, currentTime, config.AttackStartTime, config.AttackEndTime) // Each LotteryEvent represents one attempt
	if win {
		// Assign a shard based on the winning ticket
		n.AssignedShard = lottery.AssignShard(n.ID, currentTime, numShards)
		return true, n.AssignedShard
	}
	return false, -1
}

func (n *Node) IsAssignedToShard() bool {
	return n.AssignedShard != -1
}

func (n *Node) CreateBlock(previousBlockID int, currentTime int64) *block.Block {
	blkID := previousBlockID + 1
	blk := block.NewBlock(blkID, n.AssignedShard, n.ID, previousBlockID, currentTime)
	blk.IsMalicious = !n.IsHonest // Mark if block is malicious
	return blk
}

func (n *Node) BroadcastBlock(blk *block.Block, peers []*Node, currentTime int64) []*event.Event {
	events := make([]*event.Event, 0)
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			delay := utils.SimulateNetworkBlockDelay()
			e := &event.Event{
				Timestamp: float64(currentTime) + delay,
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
		if !blk.IsMalicious {
			n.KnownBlocks[blk.ID] = blk
		}
		// The shard's state is managed by the simulation
	}
}
