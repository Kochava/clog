package clog

import (
	"context"
	"sync"
	"testing"
)

func TestFieldGeneratorRace(t *testing.T) {
	// This test looks for a data race (presumably fixed) in
	// WithFieldGenerators.
	//
	// What happens is WithFieldGenerators will request the slice of
	// FieldGenerators from a context and append to that slice. When the
	// slice has more capacity than length, causing WithFieldGenerators to
	// modify the backing slice's array instead of allocating a new array to
	// contain the appended items.
	//
	// When this happens in parallel across multiple goroutines,
	// WithFieldGenerators will write to the same backing array, creating
	// a data race between the goroutines writing to the array. This depends
	// on WithFieldGenerators being called enough times that a backing array
	// is created with extra capacity to account for multiple appends, so
	// it's dependent on some of the internal behavior on how append selects
	// the next array's capacity.
	const goroutines = 3

	ctx := context.Background()
	ctx = WithFieldGenerators(ctx, make([]FieldGenerator, 13)...)

	var wg sync.WaitGroup
	start := make(chan struct{})

	// Simulate goroutines sharing a context over parallel tasks, such as
	// with HTTP requests, by spawning a number of them that branch off of
	// the same top-level context.
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(ictx context.Context) {
			defer wg.Done()
			<-start
			ictx = WithFieldGenerators(ictx, nil)
		}(ctx)
	}

	// Kick off the above goroutines to try to set off a data race.
	close(start)
	wg.Wait()
}
