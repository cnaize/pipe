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

var _ pipe.Pipe = (*DiscardPipe)(nil)

type DiscardPipe struct {
	*common.BasePipe
}

func Discard() *DiscardPipe {
	return &DiscardPipe{
		BasePipe: common.NewBase(),
	}
}

func (p *DiscardPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	var syncErr types.SyncError

	if state != nil {
		var wg sync.WaitGroup
		for file := range state.Files {
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
						syncErr.Join(fmt.Errorf("state: discard: copy: %w", err))
						return
					}
				}()

				return nil
			}()
			syncErr.Join(err)
		}

		state.Files = func(yield func(*types.File) bool) {}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
