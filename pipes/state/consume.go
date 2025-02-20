package state

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*ConsumePipe)(nil)

type ConsumePipe struct {
	*common.BasePipe
}

func Consume() *ConsumePipe {
	return &ConsumePipe{
		BasePipe: common.NewBase(),
	}
}

func (p *ConsumePipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	var syncErr types.SyncError

	if state != nil {
		var wg sync.WaitGroup

		var files []*types.File
		for file := range state.Files {
			files = append(files, file)

			err := func() error {
				if file == nil || file.Data == nil {
					return nil
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					if _, err := io.Copy(io.Discard, fileData); err != nil {
						syncErr.Join(fmt.Errorf("state: consume: copy: %w", err))
						return
					}
				}()

				return nil
			}()
			syncErr.Join(err)
		}

		state.Files = func(yield func(*types.File) bool) {
			for _, file := range files {
				if !yield(file) {
					break
				}
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
