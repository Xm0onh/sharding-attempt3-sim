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
		ID:               id,
		IsHonest:         true,
		IsOperator:       isOperator,
		AssignedShard:    -1, // Unassigned initially
		Resources:        1,
		Blockchain:       make(map[int]*block.Block),
		BlockHeaderChain: make(map[int]*block.BlockHeader),
	}
	n.BlockHeaderChain[0] = &block.BlockHeader{ID: 0}
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

func (n *Node) BroadcastBlock(blk *block.Block, peers []*Node, currentTime int64) ([]*event.Event, float64) {
	events := make([]*event.Event, 0)
	delay := 0.0
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			delay += utils.SimulateNetworkBlockDelay()
			e := &event.Event{
				Timestamp: float64(currentTime),
				Type:      event.MessageEvent,
				NodeID:    peerNode.ID,
				Data:      blk,
			}
			events = append(events, e)
		}
	}
	return events, delay
}

func (n *Node) BroadcastBlockHeader(blk *block.BlockHeader, peers []*Node, currentTime int64) ([]*event.Event, float64) {
	events := make([]*event.Event, 0)
	delay := 0.0
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			peerNode.HandleBlockHeader(blk)
			delay += utils.SimulateNetworkBlockHeaderDelay()
			e := &event.Event{
				Timestamp: float64(currentTime) + delay,
				Type:      event.MessageEvent,
				NodeID:    peerNode.ID,
				Data:      blk,
			}
			events = append(events, e)
		}
	}
	return events, delay
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
	// fmt.Println("Handling block header for node", n.ID, "Block ID", blk.ID, "from", blk.ProducerID)
	if _, exists := n.BlockHeaderChain[blk.ID]; !exists {
		n.BlockHeaderChain[blk.ID] = blk
	}
}

func (n *Node) LatestBlockHeaderID() int {
	if len(n.BlockHeaderChain) == 0 {
		n.BlockHeaderChain[0] = &block.BlockHeader{ID: 0}
		return 0
	}
	// Return the ID of the last block header
	return n.BlockHeaderChain[len(n.BlockHeaderChain)-1].ID
}

func (n *Node) DownloadLatestKBlocks(peers []*Node, currentTime int64) float64 {
	// Get our latest block header ID
	latestID := n.LatestBlockHeaderID()

	// Calculate the range of blocks we need to download
	startID := max(0, latestID-config.NumBlocksToDownload)

	totalDelay := 0.0
	downloadedBlocks := make(map[int]bool)

	// Try to download blocks from peers
	for blockID := latestID; blockID > startID; blockID-- {
		if _, exists := n.Blockchain[blockID]; exists {
			continue // Skip if we already have this block
		}

		for _, peer := range peers {
			if block, exists := peer.Blockchain[blockID]; exists {
				if !downloadedBlocks[blockID] {
					delay := utils.SimulateNetworkBlockDownloadDelay()
					if !peer.IsHonest {
						delay += float64(config.TimeOut)
					}
					totalDelay += delay

					// Only store honest blocks
					if !block.IsMalicious {
						n.Blockchain[blockID] = block
					}

					downloadedBlocks[blockID] = true
					break // Move to next block once we've downloaded this one
				}
			}
		}
	}

	return totalDelay
}
