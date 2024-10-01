package event

import (
	"container/heap"
)

type EventType int

const (
	LotteryEvent EventType = iota
	MessageEvent
	AttackEvent
	MetricsEvent
)

type Event struct {
	Timestamp int64
	Type      EventType
	NodeID    int
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
	// Min-heap based on Timestamp
	return eq[i].Timestamp < eq[j].Timestamp
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *EventQueue) Push(x interface{}) {
	// Push uses pointer receiver because it modifies the slice
	*eq = append(*eq, x.(*Event))
}

func (eq *EventQueue) Pop() interface{} {
	// Pop removes and returns the last element
	old := *eq
	n := len(old)
	x := old[n-1]
	*eq = old[0 : n-1]
	return x
}

func (eq *EventQueue) IsEmpty() bool {
	return eq.Len() == 0
}
