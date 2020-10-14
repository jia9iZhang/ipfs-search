package worker

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

// Group represents a group of Count workers, created by Factory
type Group struct {
	Count   uint
	Factory Factory
	Wait    time.Duration // Time to wait between starting workers
}

// Work starts Count of workers, created by Factory
func (g *Group) Work(ctx context.Context) error {
	var (
		worker Worker
		err    error
	)

	// Create error group and context
	errg, ctx := errgroup.WithContext(ctx)

	// Create a pool of workers within errorgroup
	for i := uint(0); i < g.Count; i++ {
		log.Printf("Initialising worker %s (%d)", worker, i+1)
		worker, err = g.Factory()
		if err != nil {
			return err
		}

		log.Printf("Starting worker %s (%d)", worker, i+1)
		errg.Go(func() error {
			return worker.Work(ctx)
		})

		time.Sleep(g.Wait)
	}

	// Block until done, returning an error if and as soon as one of the
	// child contexts errors
	return errg.Wait()
}
