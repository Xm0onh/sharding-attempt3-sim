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
---> DOWNLOAD LATEST K BLOCKS
---> Broadcast block body to operators within the shard ONLY
*/
type Simulation struct {
	Config                             config.Config
	Nodes                              map[int]*node.Node
	Operators                          map[int]*node.Node
	Shards                             map[int]*shard.Shard
	EventQueue                         *event.EventQueue
	Metrics                            *metrics.MetricsCollector
	CurrentTime                        int64
	NetworkBlockBroadcastDelays        map[int][]int64
	NetworkBlockHeaderDelays           map[int][]int64
	NetworkBlockDownloadDelays         map[int][]int64
	Logs                               []string
	currentStepMaliciousShardRotations int
	TotalRotations                     int
	NextBlockProducer                  map[int]map[int]bool
}

func NewSimulation(cfg config.Config) *Simulation {
	sim := &Simulation{
		Config:                      cfg,
		Nodes:                       make(map[int]*node.Node),
		Operators:                   make(map[int]*node.Node),
		Shards:                      make(map[int]*shard.Shard),
		EventQueue:                  event.NewEventQueue(),
		Metrics:                     metrics.NewMetricsCollector(),
		CurrentTime:                 0,
		NetworkBlockBroadcastDelays: make(map[int][]int64),
		NetworkBlockHeaderDelays:    make(map[int][]int64),
		NetworkBlockDownloadDelays:  make(map[int][]int64),
		Logs:                        make([]string, 0),
		NextBlockProducer:           make(map[int]map[int]bool),
	}

	sim.initializeNodes()
	sim.initializeOperators()
	sim.initializeShards()
	sim.initializeOperatorsMap()
	sim.scheduleInitialEvents()

	return sim
}

func (sim *Simulation) initializeNodes() {
	for i := 0; i < sim.Config.NumNodes; i++ {
		n := node.NewNode(i, false)
		sim.Nodes[n.ID] = n
	}
}

func (sim *Simulation) initializeOperators() {
	// Calculate operators per shard to ensure equal distribution
	operatorsPerShard := sim.Config.NumOperators / sim.Config.NumShards

	// Create operators and assign them to shards sequentially
	operatorID := 0
	for shardID := 0; shardID < sim.Config.NumShards; shardID++ {
		for i := 0; i < operatorsPerShard; i++ {
			n := node.NewNode(operatorID, true)
			sim.Operators[n.ID] = n
			operatorID++
		}
	}
}

func (sim *Simulation) initializeOperatorsMap() {
	fmt.Println("Initializing operators map")
	operatorsPerShard := sim.Config.NumOperators / sim.Config.NumShards
	// Assign operators to shards in groups
	operatorID := 0
	for shardID := 0; shardID < sim.Config.NumShards; shardID++ {
		for i := 0; i < operatorsPerShard; i++ {
			n := sim.Operators[operatorID]
			sim.Shards[shardID].AddNode(n)
			n.AssignedShard = shardID
			operatorID++
		}
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
			Timestamp: float64(sim.CurrentTime) + float64(sim.Config.BlockProductionInterval),
			Type:      event.ShardBlockProductionEvent,
			ShardID:   s.ID,
		}
		heap.Push(sim.EventQueue, e)
	}
	// Schedule the first LotteryEvent for all nodes
	e := &event.Event{
		Timestamp: float64(sim.CurrentTime) + float64(sim.Config.BlockProductionInterval),
		Type:      event.LotteryEvent,
	}
	heap.Push(sim.EventQueue, e)

	// Schedule the first MetricsEvent
	e = &event.Event{
		Timestamp: float64(sim.CurrentTime),
		Type:      event.MetricsEvent,
	}
	heap.Push(sim.EventQueue, e)
}

func (sim *Simulation) Run() {
	for !sim.EventQueue.IsEmpty() && sim.CurrentTime <= (sim.Config.SimulationTime+sim.Config.BlockProductionInterval) {
		e := heap.Pop(sim.EventQueue).(*event.Event)
		sim.CurrentTime = int64(e.Timestamp)
		sim.processEvent(e)
	}
	// Network delay for each shard

	fmt.Println("Network delay for shard ->", sim.NetworkBlockBroadcastDelays)
	sim.handleMetricsEvent()

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
		Timestamp: float64(sim.CurrentTime) + float64(sim.Config.BlockProductionInterval),
		Type:      event.LotteryEvent,
	}
	heap.Push(sim.EventQueue, nextEvent)
}

func (sim *Simulation) processLotteryWin(n *node.Node, newShardID int) {
	oldShardID := n.AssignedShard
	log := fmt.Sprintf("[Lottery] Node %d won the lottery and moved from Shard %d to Shard %d at time %d", n.ID, oldShardID, newShardID, sim.CurrentTime)
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
	sim.NextBlockProducer[newShardID][n.ID] = false

}

func (sim *Simulation) handleShardBlockProductionEvent(e *event.Event) {
	shardID := e.ShardID
	s := sim.Shards[shardID]

	if len(sim.NextBlockProducer[shardID]) == 0 {
		// No nodes assigned to this shard
		// Schedule next ShardBlockProductionEvent
		nextEvent := &event.Event{
			Timestamp: float64(sim.CurrentTime) + float64(sim.Config.BlockProductionInterval),
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
		latestBlockID := s.LatestBlockID()
		blk := producerNode.CreateBlock(latestBlockID, sim.CurrentTime)
		blkHeader := producerNode.CreateBlockHeader(latestBlockID, sim.CurrentTime)
		// The proposer must add the block to its blockchain
		producerNode.HandleBlock(blk)

		blk.Timestamp = sim.CurrentTime

		// Node broadcasts the block to peers in the shard
		shardOperatorNodes := sim.getShardOperators(shardID)
		events := producerNode.BroadcastBlock(blk, shardOperatorNodes, sim.CurrentTime)

		var totalDelay float64
		for _, evt := range events {
			heap.Push(sim.EventQueue, evt)
			// AVERAGE DELAY
			delay := evt.Timestamp - float64(sim.CurrentTime)
			totalDelay += delay
		}
		if len(events) > 0 {
			sim.NetworkBlockBroadcastDelays[shardID] = append(sim.NetworkBlockBroadcastDelays[shardID], int64(totalDelay/float64(len(events))))
		}

		log := fmt.Sprintf("[Block Production] Node %d produced block %d at time %d", producerNode.ID, blk.ID, sim.CurrentTime)
		sim.Logs = append(sim.Logs, log)
		sim.NextBlockProducer[shardID][producerNode.ID] = true

		// Broadcast block header to all nodes in the whole network
		events = producerNode.BroadcastBlockHeader(blkHeader, sim.getNodes(), sim.CurrentTime)
		totalDelay = 0
		for _, evt := range events {
			delay := evt.Timestamp - float64(sim.CurrentTime)
			totalDelay += delay
		}
		if len(events) > 0 {
			sim.NetworkBlockHeaderDelays[shardID] = append(sim.NetworkBlockHeaderDelays[shardID], int64(totalDelay/float64(len(events))))
		}

		// Schedule next ShardBlockProductionEvent for this shard
		nextEvent := &event.Event{
			Timestamp: float64(sim.CurrentTime) + float64(sim.Config.BlockProductionInterval),
			Type:      event.ShardBlockProductionEvent,
			ShardID:   shardID,
		}
		heap.Push(sim.EventQueue, nextEvent)
	}

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
	// sim.Metrics.Collect(sim.CurrentTime, sim.Shards, sim.Nodes, sim.NetworkBlockBroadcastDelays, sim.Logs, sim.currentStepMaliciousShardRotations)
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

func (sim *Simulation) getNodes() []*node.Node {
	nodes := []*node.Node{}
	for _, n := range sim.Nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (sim *Simulation) getShardOperators(shardID int) []*node.Node {
	nodes := []*node.Node{}
	for _, n := range sim.Operators {
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
