package internals

import "sync"

// GenericWorkers spins up given number of genericWorker routines
// that all perform work and it itself returns the outputs channel
// with the buffer defined by cz.
func GenericWorkers[A, B any](
	inputs <-chan A,
	work func(A) B,
	workers int,
	cz int,
) <-chan B {
	outputs := make(chan B, cz)
	wg := &sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go GenericWorker(inputs, outputs, work, wg)
	}
	go func() {
		wg.Wait()
		close(outputs)
	}()
	return outputs
}

// GenericWorker gets values from inputs channel, performs work
// on the value and puts it into the outputs channel, while also
// marking itself as Done in WaitGroup when channel is closed.
func GenericWorker[A, B any](
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
