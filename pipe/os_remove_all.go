package pipe

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipes/types"
)

type OsRemoveAllPipe struct {
	*basePipe

	path string
}

func OsRemoveAll(path string) *OsRemoveAllPipe {
	return &OsRemoveAllPipe{
		basePipe: newBasePipe(),
		path:     path,
	}
}

func (p *OsRemoveAllPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	if err := os.RemoveAll(p.path); err != nil {
		return nil, fmt.Errorf("os remove all: remove all: %w", err)
	}

	return p.basePipe.Send(ctx, in)
}
