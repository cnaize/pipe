package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*WriteToPipe)(nil)

type WriteToPipe struct {
	*BasePipe

	writers []io.Writer
}

func WriteTo(writers ...io.Writer) *WriteToPipe {
	return &WriteToPipe{
		BasePipe: NewBase(),
		writers:  writers,
	}
}

func (p *WriteToPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var mu sync.Mutex
	var errs error

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		next, stop := iter.Pull2(files)
		defer stop()

		var wg sync.WaitGroup
		for _, writer := range p.writers {
			if ok := func() bool {
				file, err, ok := next()
				if err != nil {
					return yield(nil, fmt.Errorf("common: write to: next: %w", err))
				}
				if !ok {
					return false
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					n, err := io.Copy(writer, fileData)
					if err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("localfs: create files: copy: %w", err))
						return
					}

					file.Size = n
				}()

				return yield(file, nil)
			}(); !ok {
				break
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
