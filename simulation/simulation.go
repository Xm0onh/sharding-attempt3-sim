// simulation.go

package simulation

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sharding/attack"
	"sharding/block"
	"sharding/config"
	"sharding/event"
	"sharding/metrics"
	"sharding/node"
	"sharding/shard"
)

type Simulation struct {
	Config                             config.Config
	Nodes                              map[int]*node.Node
	Shards                             map[int]*shard.Shard
	EventQueue                         *event.EventQueue
	Metrics                            *metrics.MetricsCollector
	CurrentTime                        int64
	NetworkDelays                      []int64 // For network latency statistics
	AttackLogs                         []string
	currentStepMaliciousShardRotations int
	TotalRotations                     int
}

func NewSimulation(cfg config.Config) *Simulation {
	sim := &Simulation{
		Config:        cfg,
		Nodes:         make(map[int]*node.Node),
		Shards:        make(map[int]*shard.Shard),
		EventQueue:    event.NewEventQueue(),
		Metrics:       metrics.NewMetricsCollector(),
		CurrentTime:   0,
		NetworkDelays: make([]int64, 0, 1000),
		AttackLogs:    make([]string, 0),
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
	}
}

func (sim *Simulation) scheduleInitialEvents() {
	// Schedule initial ShardBlockProductionEvents for each shard
	for _, s := range sim.Shards {
		e := &event.Event{
			Timestamp: sim.CurrentTime + rand.Int63n(sim.Config.BlockProductionInterval), // Stagger initial block production
			Type:      event.ShardBlockProductionEvent,
			ShardID:   s.ID,
		}
		heap.Push(sim.EventQueue, e)
	}

	// Schedule initial LotteryEvents for all nodes
	for _, n := range sim.Nodes {
		e := &event.Event{
			Timestamp: sim.CurrentTime + rand.Int63n(sim.Config.TimeStep+1), // Stagger initial lottery
			Type:      event.LotteryEvent,
			NodeID:    n.ID,
		}
		heap.Push(sim.EventQueue, e)
	}

	// Schedule the first MetricsEvent
	e := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.MetricsEvent,
	}
	heap.Push(sim.EventQueue, e)

	// Schedule AttackEvents based on AttackSchedule
	for atkTime, atkType := range sim.Config.AttackSchedule {
		e := &event.Event{
			Timestamp: atkTime,
			Type:      event.AttackEvent,
			Data:      atkType, // Pass the attack type as data
		}
		heap.Push(sim.EventQueue, e)
	}
}

func (sim *Simulation) Run() {
	for !sim.EventQueue.IsEmpty() && sim.CurrentTime <= sim.Config.SimulationTime {
		e := heap.Pop(sim.EventQueue).(*event.Event) // Type assertion added
		sim.CurrentTime = e.Timestamp
		sim.processEvent(e)
	}
}

func (sim *Simulation) processEvent(e *event.Event) {
	switch e.Type {
	case event.LotteryEvent:
		sim.handleLotteryEvent(e)
	case event.ShardBlockProductionEvent:
		sim.handleShardBlockProductionEvent(e)
	case event.MessageEvent:
		sim.handleMessageEvent(e)
	case event.MetricsEvent:
		sim.handleMetricsEvent()
	default:
		// Unknown event type
		log := fmt.Sprintf("[Simulation] Unknown event type at time %d", sim.CurrentTime)
		sim.AttackLogs = append(sim.AttackLogs, log)
	}
}

func (sim *Simulation) handleLotteryEvent(e *event.Event) {
	n := sim.Nodes[e.NodeID]
	won, newShardID := n.ParticipateInLottery(sim.CurrentTime, sim.Config.NumShards)
	maliciousShardRotation := 0

	if won {
		oldShardID := n.AssignedShard
		log := fmt.Sprintf("[Simulation] Node %d won the lottery and moved from Shard %d to Shard %d at time %d", n.ID, oldShardID, newShardID, sim.CurrentTime)
		sim.AttackLogs = append(sim.AttackLogs, log)

		sim.TotalRotations++

		if !n.IsHonest {
			maliciousShardRotation = 1
		}

		// Remove node from old shard if it was assigned to one
		if oldShardID != -1 {
			oldShard := sim.Shards[oldShardID]
			oldShard.RemoveNode(n.ID)
		}

		// Assign node to the new shard
		newShard := sim.Shards[newShardID]
		newShard.AddNode(n)
		n.AssignedShard = newShardID

		// Node produces a block immediately upon assignment to new shard
		latestBlockID := newShard.LatestBlockID()
		blk := n.CreateBlock(latestBlockID, sim.CurrentTime)
		newShard.AddBlock(blk)

		// Node broadcasts the block to peers in the new shard
		shardNodes := sim.getShardNodes(newShardID)
		events := n.BroadcastBlock(blk, shardNodes, sim.CurrentTime)
		for _, evt := range events {
			heap.Push(sim.EventQueue, evt)
			sim.NetworkDelays = append(sim.NetworkDelays, evt.Timestamp-sim.CurrentTime)
		}
	}

	// Schedule the next LotteryEvent for this node
	nextEvent := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.LotteryEvent,
		NodeID:    n.ID,
	}
	heap.Push(sim.EventQueue, nextEvent)

	sim.currentStepMaliciousShardRotations += maliciousShardRotation
}

func (sim *Simulation) handleShardBlockProductionEvent(e *event.Event) {
	shardID := e.ShardID
	s := sim.Shards[shardID]
	shardNodes := sim.getShardNodes(shardID)

	if len(shardNodes) == 0 {
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

	// Randomly select a node from shardNodes
	selectedNode := shardNodes[rand.Intn(len(shardNodes))]

	// Node creates a block
	latestBlockID := s.LatestBlockID()
	blk := selectedNode.CreateBlock(latestBlockID, sim.CurrentTime)
	blk.Timestamp = sim.CurrentTime // Ensure block timestamp is set
	s.AddBlock(blk)

	// Node broadcasts the block to peers in the shard
	shardNodes = sim.getShardNodes(shardID) // Refresh shard nodes after potential changes
	events := selectedNode.BroadcastBlock(blk, shardNodes, sim.CurrentTime)
	for _, evt := range events {
		heap.Push(sim.EventQueue, evt)
		sim.NetworkDelays = append(sim.NetworkDelays, evt.Timestamp-sim.CurrentTime)
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

func (sim *Simulation) handleAttackEvent(e *event.Event) {
	atkType, ok := e.Data.(config.AttackType)
	if !ok {
		log := fmt.Sprintf("[Simulation] Invalid attack data at time %d", sim.CurrentTime)
		sim.AttackLogs = append(sim.AttackLogs, log)
		return
	}

	attack.ExecuteAttack(atkType, sim.CurrentTime, sim.Nodes, sim.Shards, sim.EventQueue, sim.Config, &sim.AttackLogs)
}

func (sim *Simulation) handleMetricsEvent() {
	// Collect metrics with the count of malicious shard rotations and total rotations in this step
	sim.Metrics.Collect(sim.CurrentTime, sim.Shards, sim.Nodes, sim.NetworkDelays, sim.AttackLogs, sim.currentStepMaliciousShardRotations)
	// fmt.Printf("Current Time: %d, Malicious Shard Rotations: %d, Percentage of malicious: %.2f%%\n",
	// 	sim.CurrentTime,
	// 	sim.currentStepMaliciousShardRotations,
	// 	float64(sim.currentStepMaliciousShardRotations)/float64(sim.TotalRotations)*100)

	sim.NetworkDelays = sim.NetworkDelays[:0]
	sim.AttackLogs = sim.AttackLogs[:0]        // Reset attack logs after collecting
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

// Additional field to track rotations in the current step
