package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	numberOfWorkers = 10
)

func main() {
	// create a consumer
	consumer := Consumer{
		ingestChan: make(chan JobEvent, 1),
		jobsChan:   make(chan JobEvent, numberOfWorkers),
	}

	// simulate external lib sending jobs
	producer := Producer{CallbackFunc: consumer.RegisterEventCallback}
	go producer.StartSimulation()

	// context for cancellation
	ctx, cancelFunc := context.WithCancel(context.Background())

	// use a wait group to shutdown program only when all jobs are done
	wg := &sync.WaitGroup{}

	// start the consumer
	go consumer.Start(ctx)

	// add num of workers to wg
	wg.Add(numberOfWorkers)
	// start the worker pool
	for i := 0; i < numberOfWorkers; i++ {
		go consumer.Worker(wg, i)
	}

	// we want to create a termination channel so we can shut down gracefully
	terminationChan := make(chan os.Signal)
	signal.Notify(terminationChan, syscall.SIGINT, syscall.SIGTERM)

	<-terminationChan

	fmt.Println("--- SHUTDOWN SIGNAL RECEIVED ---")
	cancelFunc() // send shutdown signal through the context
	wg.Wait()    // wait here until all workers are done

	fmt.Println("All workers finished")
}
