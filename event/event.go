package event

import (
	"container/heap"
)

type EventType int

const (
	LotteryEvent EventType = iota
	ShardBlockProductionEvent
	MessageEvent
	MetricsEvent
)

type Event struct {
	Timestamp float64
	Type      EventType
	NodeID    int
	ShardID   int
	Data      interface{}
}

type EventQueue []*Event

func NewEventQueue() *EventQueue {
	eq := &EventQueue{}
	heap.Init(eq)
	return eq
}

func (eq EventQueue) Len() int { return len(eq) }

func (eq EventQueue) Less(i, j int) bool {
	if eq[i].Timestamp == eq[j].Timestamp {
		return eq[i].Type < eq[j].Type
	}
	return eq[i].Timestamp < eq[j].Timestamp
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *EventQueue) Push(x interface{}) {
	*eq = append(*eq, x.(*Event))
}

func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	item := old[0]
	old[0] = old[n-1]
	*eq = old[0 : n-1]
	return item
}

func (eq *EventQueue) IsEmpty() bool {
	return eq.Len() == 0
}
