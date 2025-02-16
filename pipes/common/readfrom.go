package common

import (
	"context"
	"io"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*ReadFromPipe)(nil)

type ReadFromPipe struct {
	*BasePipe

	readers []io.Reader
}

func ReadFrom(readers ...io.Reader) *ReadFromPipe {
	return &ReadFromPipe{
		BasePipe: NewBase(),
		readers:  readers,
	}
}

func (p *ReadFromPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		for file, err := range files {
			if !yield(file, err) {
				break
			}
		}

		for _, reader := range p.readers {
			if !yield(&types.File{Data: reader}, nil) {
				break
			}
		}
	}

	return p.BasePipe.Run(ctx, state)
}
