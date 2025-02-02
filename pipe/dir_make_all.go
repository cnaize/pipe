package pipe

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipes/types"
)

type DirMakeAllPipe struct {
	*basePipe

	path string
	perm os.FileMode
}

func DirMakeAll(path string, perm os.FileMode) *DirMakeAllPipe {
	return &DirMakeAllPipe{
		basePipe: newBasePipe(),
		path:     path,
		perm:     perm,
	}
}

func (p *DirMakeAllPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	if err := os.MkdirAll(p.path, p.perm); err != nil {
		return nil, fmt.Errorf("dir make all: mkdir all: %w", err)
	}

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	out.DirMakeAll = &types.Dir{
		Path: p.path,
		Perm: p.perm,
	}

	return out, nil
}
