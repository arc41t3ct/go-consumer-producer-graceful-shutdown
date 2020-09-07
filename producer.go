package main

import (
	"fmt"
	"time"
)

// Producer simulates an external library that invokes the registered
// callback when it has new data for us once per 100ms
type Producer struct {
	CallbackFunc func(event JobEvent)
}

func (p Producer) StartSimulation() {
	eventIndex := 0
	for {
		p.CallbackFunc(JobEvent{
			Id:   fmt.Sprintf("%d", eventIndex),
			Type: fmt.Sprintf("Simulation %v", eventIndex),
		})
		eventIndex++
		time.Sleep(time.Millisecond * 100)
	}
}
