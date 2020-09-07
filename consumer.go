package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type JobEvent struct {
	Id   string
	Type string
}

type Consumer struct {
	ingestChan chan JobEvent
	jobsChan   chan JobEvent
}

func (c Consumer) Start(ctx context.Context) {
	for {
		select {
		case job := <-c.ingestChan:
			c.jobsChan <- job
		case <-ctx.Done():
			fmt.Println("Consumer received cancellation signal, closing jobsChan")
			close(c.jobsChan)
			fmt.Println("Consumer closed jobsChan")
			return
		}
	}
}

// Externally register an event
func (c Consumer) RegisterEventCallback(event JobEvent) {
	c.ingestChan <- event
}

// starts a single worker function
func (c Consumer) Worker(wg *sync.WaitGroup, workerIndex int) {
	fmt.Printf("Worker %d starting\n", workerIndex)
	// defer will be called once the jobsChan closes
	defer wg.Done()
	// blocks until an event is received or channel is closed
	for event := range c.jobsChan {
		// handle events
		fmt.Printf("Worker %d started job %v type %v\n", workerIndex, event.Id, event.Type)
		time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
		fmt.Printf("Worker %d finished job %v type %v\n", workerIndex, event.Id, event.Type)
	}
	fmt.Printf("Worker %d interrupted\n", workerIndex)
}
