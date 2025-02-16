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

var _ pipe.Pipe = (*ConsumeFilesPipe)(nil)

type ConsumeFilesPipe struct {
	*common.BasePipe
}

func ConsumeFiles() *ConsumeFilesPipe {
	return &ConsumeFilesPipe{
		BasePipe: common.NewBase(),
	}
}

func (p *ConsumeFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	var mu sync.Mutex
	var errs error

	if state != nil {
		var wg sync.WaitGroup
		var files []*types.File
		for file, err := range state.Files {
			files = append(files, file)

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

		state.Files = func(yield func(*types.File, error) bool) {
			for _, file := range files {
				if !yield(file, nil) {
					break
				}
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
