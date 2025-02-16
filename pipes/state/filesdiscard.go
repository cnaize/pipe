package state

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*DiscardFilesPipe)(nil)

type DiscardFilesPipe struct {
	*common.BasePipe
}

func DiscardFiles() *DiscardFilesPipe {
	return &DiscardFilesPipe{
		BasePipe: common.NewBase(),
	}
}

func (p *DiscardFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	var mu sync.Mutex
	var errs error

	if state != nil {
		var wg sync.WaitGroup
		for file, err := range state.Files {
			func() {
				if err != nil {
					mu.Lock()
					defer mu.Unlock()

					errs = errors.Join(errs, fmt.Errorf("state: consume files: next: %w", err))
					return
				}
				if file == nil || file.Data == nil {
					return
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					if _, err := io.Copy(io.Discard, fileData); err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("state: consume files: copy: %w", err))
						return
					}
				}()
			}()
		}

		state.Files = func(yield func(*types.File, error) bool) {}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
