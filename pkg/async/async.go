package async

import (
	"context"
	"io"

	"log"

	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run(ctx context.Context) error
	io.Closer
}

type Func func(ctx context.Context) error

func (f Func) Run(ctx context.Context) error {
	return f(ctx)
}

func (f Func) Close() error {
	return nil
}

func RunGroup(rg []Runner) Runner {
	ru := Func(func(ctx context.Context) error {
		g, ctx := errgroup.WithContext(ctx)
		for idx := range rg {
			r := rg[idx]
			g.Go(func() error {
				defer func() {
					if err := r.Close(); err != nil {
						log.Printf("Error when closing: %v", err)
					}
				}()

				return r.Run(ctx)
			})
		}
		return g.Wait()
	})
	return ru
}
