package main

import (
	"io/fs"
	"sync"
)

const (
	// savePerms tells us what permissions to use for the
	// final export files.
	savePerms = fs.FileMode(0644)
)

var (
	// defaultNumOfWorkers gives us the number of workers to
	// spin up in each stage: parsing and processing.
	defaultNumOfWorkers = 14
	// channelCapacity is the default capacity on workers'
	// output channels.
	channelCapacity = defaultNumOfWorkers
)

// genericWorkers spins up given number of genericWorker routines
// that all perform work and it itself returns the outputs channel.
func genericWorkers[A, B any](
	inputs <-chan A,
	work func(A) B,
	workers int,
) <-chan B {
	outputs := make(chan B, channelCapacity)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go genericWorker(inputs, outputs, work, wg)
	}
	go func() {
		wg.Wait()
		close(outputs)
	}()
	return outputs
}

// genericWorker gets values from inputs channel, performs work
// on the value and puts it into the outputs channel, while also
// marking itself as Done in WaitGroup when channel is closed.
func genericWorker[A, B any](
	inputs <-chan A,
	outputs chan<- B,
	work func(A) B,
	wg *sync.WaitGroup,
) {
	for input := range inputs {
		outputs <- work(input)
	}
	wg.Done()
}
