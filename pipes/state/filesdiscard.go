package state

import (
	"context"
	"fmt"
	"io"

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
	if state != nil {
		for file, err := range state.Files {
			if err != nil {
				return nil, fmt.Errorf("state: discard fiels: %w", err)
			}
			if file == nil || file.Data == nil {
				continue
			}

			if _, err := io.Copy(io.Discard, file.Data); err != nil {
				return nil, fmt.Errorf("state: discard files: copy: %w", err)
			}

			file.Data = nil
		}

		state.Files = func(yield func(*types.File, error) bool) {}
	}

	return p.BasePipe.Run(ctx, state)
}
