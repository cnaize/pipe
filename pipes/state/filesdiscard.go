package pipe

import (
	"context"
	"fmt"
	"io"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*DiscardFilesPipe)(nil)

type DiscardFilesPipe struct {
	*general.BasePipe
}

func DiscardFiles() *DiscardFilesPipe {
	return &DiscardFilesPipe{
		BasePipe: general.NewBase(),
	}
}

func (p *DiscardFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state != nil {
		for file := range state.Files {
			if !(file != nil && file.Data != nil) {
				continue
			}

			if _, err := io.Copy(io.Discard, file.Data); err != nil {
				return nil, fmt.Errorf("state: discard files: copy: %w", err)
			}
		}
	}

	return p.BasePipe.Run(ctx, state)
}
