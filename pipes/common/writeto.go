package common

import (
	"context"
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

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		var wg sync.WaitGroup
		for _, writer := range p.writers {
			ok, err := func() (bool, error) {
				file, ok := next()
				if !ok {
					return false, nil
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					size, err := io.Copy(writer, fileData)
					if err != nil {
						syncErr.Join(fmt.Errorf("common: write to: copy: %w", err))
						return
					}

					file.Size = size
				}()

				return yield(file), nil
			}()
			syncErr.Join(err)
			if !ok {
				break
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
