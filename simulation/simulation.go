package simulation

import (
	"container/heap"
	"sharding/attack"
	"sharding/block"
	"sharding/config"
	"sharding/event"
	"sharding/metrics"
	"sharding/node"
	"sharding/shard"
)

type Simulation struct {
	Config        config.Config
	Nodes         map[int]*node.Node
	Shards        map[int]*shard.Shard
	EventQueue    *event.EventQueue
	Metrics       *metrics.MetricsCollector
	CurrentTime   int64
	NetworkDelays []int64 // For network latency statistics
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
	for _, n := range sim.Nodes {
		e := &event.Event{
			Timestamp: sim.CurrentTime,
			Type:      event.LotteryEvent,
			NodeID:    n.ID,
		}
		heap.Push(sim.EventQueue, e)
	}

	e := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.MetricsEvent,
	}
	heap.Push(sim.EventQueue, e)
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
	case event.MessageEvent:
		sim.handleMessageEvent(e)
	case event.AttackEvent:
		attack.ExecuteAttack(sim.Config.AttackType, sim.CurrentTime, sim.Nodes, sim.Shards, sim.EventQueue)
	case event.MetricsEvent:
		sim.handleMetricsEvent(e)
	default:
		// Unknown event type
	}
}

func (sim *Simulation) handleLotteryEvent(e *event.Event) {
	n := sim.Nodes[e.NodeID]
	won := n.ParticipateInLottery(sim.CurrentTime, sim.Config.NumShards)
	if won {
		s := sim.Shards[n.AssignedShard]
		latestBlockID := s.LatestBlockID()
		blk := n.CreateBlock(latestBlockID, sim.CurrentTime)
		s.AddBlock(blk)

		peers := sim.getShardNodes(n.AssignedShard)
		events := n.BroadcastBlock(blk, peers, sim.CurrentTime)
		for _, evt := range events {
			heap.Push(sim.EventQueue, evt)
			sim.NetworkDelays = append(sim.NetworkDelays, evt.Timestamp-sim.CurrentTime)
		}
	}

	nextEvent := &event.Event{
		Timestamp: sim.CurrentTime + sim.Config.TimeStep,
		Type:      event.LotteryEvent,
		NodeID:    n.ID,
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

func (sim *Simulation) handleMetricsEvent(e *event.Event) {
	sim.Metrics.Collect(sim.CurrentTime, sim.Shards, sim.Nodes, sim.NetworkDelays)
	sim.NetworkDelays = sim.NetworkDelays[:0]
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
