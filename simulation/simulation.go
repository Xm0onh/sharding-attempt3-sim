// simulation.go

package simulation

import (
	"container/heap"
	"fmt"
	"sharding/block"
	"sharding/config"
	"sharding/event"
	"sharding/metrics"
	"sharding/node"
	"sharding/shard"
)

/*
TODO -
---> Distribute the nodes evenly across shards
*/
type Simulation struct {
	Config                             config.Config
	Nodes                              map[int]*node.Node
	Shards                             map[int]*shard.Shard
	EventQueue                         *event.EventQueue
	Metrics                            *metrics.MetricsCollector
	CurrentTime                        int64
	NetworkDelays                      []int64 // For network latency statistics
	Logs                               []string
	currentStepMaliciousShardRotations int
	TotalRotations                     int
	NextBlockProducer                  map[int]map[int]bool // Map[shardID][nodeID]bool
}

func NewSimulation(cfg config.Config) *Simulation {
	sim := &Simulation{
		Config:            cfg,
		Nodes:             make(map[int]*node.Node),
		Shards:            make(map[int]*shard.Shard),
		EventQueue:        event.NewEventQueue(),
		Metrics:           metrics.NewMetricsCollector(),
		CurrentTime:       0,
		NetworkDelays:     make([]int64, 0, 1000),
		Logs:              make([]string, 0),
		NextBlockProducer: make(map[int]map[int]bool),
	}

	sim.initializeNodes()
	sim.initializeShards()
	sim.scheduleInitialEvents()

	return sim
}

func (sim *Simulation) initializeNodes() {
	for i := 0; i < sim.Config.NumNodes; i++ {
		n := node.NewNode(i)
		sim.Nodes[n.ID] = n
	}
}

func (sim *Simulation) initializeShards() {
	for i := 0; i < sim.Config.NumShards; i++ {
		s := shard.NewShard(i)
		sim.Shards[s.ID] = s
		sim.NextBlockProducer[s.ID] = make(map[int]bool)

	}
}

func (sim *Simulation) scheduleInitialEvents() {
	for _, s := range sim.Shards {
		e := &event.Event{
			Timestamp: sim.CurrentTime + sim.Config.BlockProductionInterval,
			Type:      event.ShardBlockProductionEvent,
			ShardID:   s.ID,
		}
		heap.Push(sim.EventQueue, e)
	}
	// Schedule the first LotteryEvent for all nodes
	e := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.BlockProductionInterval,
		Type:      event.LotteryEvent,
	}
	heap.Push(sim.EventQueue, e)

	// Schedule the first MetricsEvent
	e = &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.MetricsEvent,
	}
	heap.Push(sim.EventQueue, e)
}

func (sim *Simulation) Run() {
	var counter = 0
	// Print event types in the event queue
	// for !sim.EventQueue.IsEmpty() {
	// 	e := heap.Pop(sim.EventQueue).(*event.Event)
	// 	fmt.Println("Event ->", e.Type)
	// }

	for !sim.EventQueue.IsEmpty() && sim.CurrentTime <= sim.Config.SimulationTime {
		e := heap.Pop(sim.EventQueue).(*event.Event)
		sim.CurrentTime = e.Timestamp
		sim.processEvent(e)
		// 0 -> LotteryEvent
		// 1 -> MessageEvent
		// 2 -> MetricsEvent
		// 3 -> ShardBlockProductionEvent
		if e.Type == 1 {
			counter++
		}
	}
	fmt.Println("Counter ->", counter)
}

func (sim *Simulation) processEvent(e *event.Event) {
	switch e.Type {
	case event.LotteryEvent:
		sim.handleLotteryEvent()
	case event.ShardBlockProductionEvent:
		sim.handleShardBlockProductionEvent(e)
	case event.MessageEvent:
		sim.handleMessageEvent(e)
	case event.MetricsEvent:
		sim.handleMetricsEvent()
	default:
		// Unknown event type
		log := fmt.Sprintf("[Simulation] Unknown event type at time %d", sim.CurrentTime)
		sim.Logs = append(sim.Logs, log)
	}
}

func (sim *Simulation) handleLotteryEvent() {
	for _, n := range sim.Nodes {
		won, newShardID := n.ParticipateInLottery(sim.CurrentTime, sim.Config.NumShards)
		if won {
			sim.processLotteryWin(n, newShardID)
		}
	}

	// Schedule the next LotteryEvent for all nodes
	nextEvent := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.BlockProductionInterval,
		Type:      event.LotteryEvent,
	}
	heap.Push(sim.EventQueue, nextEvent)
}

func (sim *Simulation) processLotteryWin(n *node.Node, newShardID int) {
	oldShardID := n.AssignedShard
	log := fmt.Sprintf("[Simulation] Node %d won the lottery and moved from Shard %d to Shard %d at time %d", n.ID, oldShardID, newShardID, sim.CurrentTime)
	sim.Logs = append(sim.Logs, log)
	sim.TotalRotations++

	if !n.IsHonest {
		sim.currentStepMaliciousShardRotations++
	}

	// Remove node from old shard if it was assigned to one
	if oldShardID != -1 {
		oldShard := sim.Shards[oldShardID]
		oldShard.RemoveNode(n.ID)
		delete(sim.NextBlockProducer[oldShardID], n.ID)

	}

	// Assign node to the new shard
	newShard := sim.Shards[newShardID]
	newShard.AddNode(n)
	n.AssignedShard = newShardID

	// Set the next block producer for the new shard
	sim.NextBlockProducer[newShardID][n.ID] = false
	// fmt.Println("NextBlockProducer ->", sim.NextBlockProducer)
	//TODO - DOWNLOAD LATEST K BLOCKS
	//
	//
	//
	//
	//

	// Node produces a block immediately upon assignment to new shard
	latestBlockID := newShard.LatestBlockID()
	blk := n.CreateBlock(latestBlockID, sim.CurrentTime)
	newShard.AddBlock(blk)

	/*
		TODO -
		- Broadcast block header to all the nodes
		- Broadcast block body to operators within the shard ONLY
	*/
	// Node broadcasts the block to peers in the new shard
	shardNodes := sim.getShardNodes(newShardID)
	events := n.BroadcastBlock(blk, shardNodes, sim.CurrentTime)
	for _, evt := range events {
		heap.Push(sim.EventQueue, evt)
		sim.NetworkDelays = append(sim.NetworkDelays, evt.Timestamp-sim.CurrentTime)
	}
}
func (sim *Simulation) handleShardBlockProductionEvent(e *event.Event) {
	shardID := e.ShardID
	s := sim.Shards[shardID]

	if len(sim.NextBlockProducer[shardID]) == 0 {
		// No nodes assigned to this shard
		// Schedule next ShardBlockProductionEvent
		nextEvent := &event.Event{
			Timestamp: sim.CurrentTime + sim.Config.BlockProductionInterval,
			Type:      event.ShardBlockProductionEvent,
			ShardID:   shardID,
		}
		heap.Push(sim.EventQueue, nextEvent)
		return
	}

	// Find the first node with bool == false
	var producerNode *node.Node
	for nodeID, hasProduced := range sim.NextBlockProducer[shardID] {
		if !hasProduced {
			producerNode = sim.Nodes[nodeID]
			break
		}
	}

	if producerNode == nil {
		// All nodes have produced blocks, skip producing a block
		log := fmt.Sprintf("All nodes in shard %d have produced blocks, skipping block production at time %d", shardID, sim.CurrentTime)
		sim.Logs = append(sim.Logs, log)
	} else {
		// fmt.Println("producerID ->", producerNode.ID)
		// Node creates a block
		latestBlockID := s.LatestBlockID()
		blk := producerNode.CreateBlock(latestBlockID, sim.CurrentTime)
		blk.Timestamp = sim.CurrentTime // Ensure block timestamp is set

		// Node broadcasts the block to peers in the shard
		shardNodes := sim.getShardNodes(shardID) // Get updated shard nodes
		events := producerNode.BroadcastBlock(blk, shardNodes, sim.CurrentTime)
		for _, evt := range events {
			heap.Push(sim.EventQueue, evt)
			sim.NetworkDelays = append(sim.NetworkDelays, evt.Timestamp-sim.CurrentTime)
		}

		// Set the node's bool to true
		sim.NextBlockProducer[shardID][producerNode.ID] = true
	}

	// Schedule next ShardBlockProductionEvent for this shard
	nextEvent := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.BlockProductionInterval,
		Type:      event.ShardBlockProductionEvent,
		ShardID:   shardID,
	}
	heap.Push(sim.EventQueue, nextEvent)
}

func (sim *Simulation) handleMessageEvent(e *event.Event) {
	n := sim.Nodes[e.NodeID]
	n.ProcessMessage(e)

	blk, ok := e.Data.(*block.Block)
	if ok {
		s := sim.Shards[blk.ShardID]
		s.AddBlock(blk)
	}
}

func (sim *Simulation) handleMetricsEvent() {
	// Collect metrics with the count of malicious shard rotations and total rotations in this step
	sim.Metrics.Collect(sim.CurrentTime, sim.Shards, sim.Nodes, sim.NetworkDelays, sim.Logs, sim.currentStepMaliciousShardRotations)
	// fmt.Printf("Current Time: %d, Malicious Shard Rotations: %d, Percentage of malicious: %.2f%%\n",
	// 	sim.CurrentTime,
	// 	sim.currentStepMaliciousShardRotations,
	// 	float64(sim.currentStepMaliciousShardRotations)/float64(sim.TotalRotations)*100)

	sim.NetworkDelays = sim.NetworkDelays[:0]
	sim.Logs = sim.Logs[:0]                    // Reset attack logs after collecting
	sim.currentStepMaliciousShardRotations = 0 // Reset malicious rotations count
	sim.TotalRotations = 0                     // Reset total rotations count

	// Schedule the next MetricsEvent
	nextEvent := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.MetricsEvent,
	}
	heap.Push(sim.EventQueue, nextEvent)
}

func (sim *Simulation) getShardNodes(shardID int) []*node.Node {
	nodes := []*node.Node{}
	for _, n := range sim.Nodes {
		if n.AssignedShard == shardID {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

/// FURUTE IMPROVEMENTS
// func (sim *Simulation) handleAttackEvent(e *event.Event) {
// 	atkType, ok := e.Data.(config.AttackType)
// 	if !ok {
// 		log := fmt.Sprintf("[Simulation] Invalid attack data at time %d", sim.CurrentTime)
// 		sim.Logs = append(sim.Logs, log)
// 		return
// 	}

// 	attack.ExecuteAttack(atkType, sim.CurrentTime, sim.Nodes, sim.Shards, sim.EventQueue, sim.Config, &sim.Logs)
// }
