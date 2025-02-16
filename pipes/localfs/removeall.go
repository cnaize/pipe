package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*RemoveAllPipe)(nil)

type RemoveAllPipe struct {
	*common.BasePipe

	path string
}

func RemoveAll(path string) *RemoveAllPipe {
	return &RemoveAllPipe{
		BasePipe: common.NewBase(),
		path:     path,
	}
}

func (p *RemoveAllPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if err := os.RemoveAll(p.path); err != nil {
		return nil, fmt.Errorf("localfs: os remove all: remove all: %w", err)
	}

	return p.BasePipe.Run(ctx, state)
}
