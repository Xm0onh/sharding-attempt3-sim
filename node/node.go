package node

import (
	"math/rand"
	"sharding/block"
	"sharding/config"
	"sharding/event"
	"sharding/lottery"
	"sharding/utils"
	"sync"
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

func NewNode(cfg *config.Config, id int, isOperator bool) *Node {
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
	if rand.Float64() < cfg.MaliciousNodeRatio {
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

func (n *Node) BroadcastBlock(cfg *config.Config, blk *block.Block, peers []*Node, currentTime int64) ([]*event.Event, float64) {
	events := make([]*event.Event, 0)
	delay := 0.0
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			delay += utils.SimulateNetworkBlockDelay(cfg, len(peers))
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

func (n *Node) BroadcastBlockHeader(cfg *config.Config, blk *block.BlockHeader, peers []*Node, currentTime int64) ([]*event.Event, float64) {
	events := make([]*event.Event, 0)
	delay := 0.0
	for _, peerNode := range peers {
		if peerNode.ID != n.ID {
			peerNode.HandleBlockHeader(blk)
			delay += utils.SimulateNetworkBlockHeaderDelay(cfg)
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

func (n *Node) DownloadLatestKBlocks(cfg *config.Config, peers []*Node, currentTime int64) float64 {
	latestID := n.LatestBlockHeaderID()
	startID := max(0, latestID-cfg.NumBlocksToDownload)
	counter := 0
	type downloadResult struct {
		blockID int
		block   *block.Block
		delay   float64
	}

	// Split peers into operators and regular nodes
	operators := make([]*Node, 0)
	regularPeers := make([]*Node, 0)
	for _, peer := range peers {
		if peer.IsOperator {
			operators = append(operators, peer)
		} else {
			regularPeers = append(regularPeers, peer)
		}
	}

	resultChan := make(chan downloadResult, cfg.MaxP2PConnections)
	var mu sync.Mutex
	downloadedBlocks := make(map[int]bool)
	totalDelay := 0.0

	// Process blocks in batches of size MaxP2PConnections
	for batchStart := latestID; batchStart > startID; batchStart -= cfg.MaxP2PConnections {
		counter++
		batchEnd := max(startID, batchStart-cfg.MaxP2PConnections)
		activeDLs := 0
		batchMaxDelay := 0.0

		// Start downloads for this batch
		for blockID := batchStart; blockID > batchEnd; blockID-- {
			if _, exists := n.Blockchain[blockID]; exists {
				continue
			}

			activeDLs++
			go func(bid int) {
				result := downloadResult{blockID: bid, delay: -1}

				// Try operators first
				for _, peer := range operators {
					mu.Lock()
					if downloadedBlocks[bid] {
						mu.Unlock()
						resultChan <- result
						return
					}
					mu.Unlock()

					if block, exists := peer.Blockchain[bid]; exists {
						delay := utils.SimulateNetworkBlockDownloadDelay(cfg)
						if !peer.IsHonest {
							delay += float64(cfg.TimeOut)
						}
						result.block = block
						result.delay = delay
						break
					}
				}

				// If block not found with operators, try regular peers
				if result.delay == -1 {
					for _, peer := range regularPeers {
						mu.Lock()
						if downloadedBlocks[bid] {
							mu.Unlock()
							resultChan <- result
							return
						}
						mu.Unlock()

						if block, exists := peer.Blockchain[bid]; exists {
							delay := utils.SimulateNetworkBlockDownloadDelay(cfg)
							if !peer.IsHonest {
								delay += float64(cfg.TimeOut)
							}
							result.block = block
							result.delay = delay
							break
						}
					}
				}
				resultChan <- result
			}(blockID)
		}

		// Wait for all downloads in this batch to complete
		for i := 0; i < activeDLs; i++ {
			result := <-resultChan
			if result.delay > 0 {
				mu.Lock()
				downloadedBlocks[result.blockID] = true
				if !result.block.IsMalicious {
					n.Blockchain[result.blockID] = result.block
				}
				batchMaxDelay = max(batchMaxDelay, result.delay)
				mu.Unlock()
			}
		}

		totalDelay += batchMaxDelay
	}
	return totalDelay
}
