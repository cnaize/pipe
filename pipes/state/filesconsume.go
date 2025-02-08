package state

import (
	"context"
	"fmt"
	"io"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*ConsumeFilesPipe)(nil)

type ConsumeFilesPipe struct {
	*general.BasePipe
}

func ConsumeFiles() *ConsumeFilesPipe {
	return &ConsumeFilesPipe{
		BasePipe: general.NewBase(),
	}
}

func (p *ConsumeFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state != nil {
		var files []*types.File
		for file, err := range state.Files {
			files = append(files, file)

			if err != nil {
				return nil, fmt.Errorf("state: consume files: %w", err)
			}
			if file == nil || file.Data == nil {
				continue
			}

			if _, err := io.Copy(io.Discard, file.Data); err != nil {
				return nil, fmt.Errorf("state: consume files: copy: %w", err)
			}

			file.Data = nil
		}

		state.Files = func(yield func(*types.File, error) bool) {
			for _, file := range files {
				if !yield(file, nil) {
					return
				}
			}
		}
	}

	return p.BasePipe.Run(ctx, state)
}
