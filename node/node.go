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
	ID               int
	IsHonest         bool
	IsOperator       bool
	AssignedShard    int
	Resources        int
	Blockchain       map[int]*block.Block
	BlockHeaderChain map[int]*block.BlockHeader
}

func NewNode(id int, isOperator bool) *Node {
	n := &Node{
		ID:            id,
		IsHonest:      true,
		IsOperator:    isOperator,
		AssignedShard: -1, // Unassigned initially
		Resources:     1,
		Blockchain:    make(map[int]*block.Block),
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

func (n *Node) CreateBlockHeader(previousBlockID int, currentTime int64) *block.BlockHeader {
	blkID := previousBlockID + 1
	blkHeader := block.NewBlockHeader(blkID, n.AssignedShard, n.ID, previousBlockID, currentTime)
	return blkHeader
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

func (n *Node) BroadcastBlockHeader(blk *block.BlockHeader, peers []*Node, currentTime int64) []*event.Event {
	events := make([]*event.Event, 0)
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			delay := utils.SimulateNetworkBlockHeaderDelay()
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

// Simulate downloading k blocks from multiple peers
func (n *Node) DownloadBlocks(k int, peers []*Node, currentTime int64) float64 {
	blocks := make([]*block.Block, 0)
	downloadedIDs := make(map[int]bool)
	totalDelay := 0.0
	// Try to download k blocks from peers
	for _, peerNode := range peers {
		// Skip if we're trying to download from ourselves
		if peerNode.ID == n.ID {
			continue
		}

		// Look through peer's blockchain
		for blockID, block := range peerNode.Blockchain {
			// Skip if we already have this block or if we've already downloaded it
			if _, exists := n.Blockchain[blockID]; exists {
				continue
			}
			if downloadedIDs[blockID] {
				continue
			}

			// Simulate network delay for downloading the block
			delay := utils.SimulateNetworkBlockDownloadDelay()
			totalDelay += delay
			// Add block to our downloaded list
			if !peerNode.IsHonest {
				blocks = append(blocks, block)
				totalDelay += float64(config.TimeOut)
			}

			downloadedIDs[blockID] = true

			// Break if we've downloaded enough blocks
			if len(blocks) >= k {
				return totalDelay
			}
		}
	}

	return totalDelay
}

func (n *Node) ProcessMessage(e *event.Event) {
	switch msg := e.Data.(type) {
	case *block.Block:
		n.HandleBlock(msg)
	case *block.BlockHeader:
		n.HandleBlockHeader(msg)
	default:
		// Handle other message types if necessary
	}
}

func (n *Node) HandleBlock(blk *block.Block) {
	if _, exists := n.Blockchain[blk.ID]; !exists {
		if !blk.IsMalicious {
			n.Blockchain[blk.ID] = blk
		}
		// The shard's state is managed by the simulation
	}
}

func (n *Node) HandleBlockHeader(blk *block.BlockHeader) {
	if _, exists := n.BlockHeaderChain[blk.ID]; !exists {
		n.BlockHeaderChain[blk.ID] = blk
	}
}
