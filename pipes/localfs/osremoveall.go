package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*OSRemoveAllPipe)(nil)

type OSRemoveAllPipe struct {
	*general.BasePipe

	path string
}

func OSRemoveAll(path string) *OSRemoveAllPipe {
	return &OSRemoveAllPipe{
		BasePipe: general.NewBase(),
		path:     path,
	}
}

func (p *OSRemoveAllPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if err := os.RemoveAll(p.path); err != nil {
		return nil, fmt.Errorf("lfs: os remove all: remove all: %w", err)
	}

	return p.BasePipe.Run(ctx, state)
}
