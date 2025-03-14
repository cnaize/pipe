package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*MakeDirAllPipe)(nil)

type MakeDirAllPipe struct {
	*common.BasePipe

	path string
	perm os.FileMode
}

func MakeDirAll(path string, perm os.FileMode) *MakeDirAllPipe {
	return &MakeDirAllPipe{
		BasePipe: common.NewBase(),
		path:     path,
		perm:     perm,
	}
}

func (p *MakeDirAllPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if err := os.MkdirAll(p.path, p.perm); err != nil {
		return nil, fmt.Errorf("localfs: make dir all: mkdir all: %w", err)
	}

	return p.BasePipe.Run(ctx, state)
}
