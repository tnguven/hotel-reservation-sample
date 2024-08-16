package utils

import (
	"context"
	"sync"
)

func Parallel[E any](events []E) func(func(E) bool) {
	return func(yield func(E) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(len(events))

		for _, event := range events {
			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
					if !yield(event) {
						cancel()
					}
				}

			}()
		}

		wg.Wait()
	}
}
