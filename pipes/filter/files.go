package filter

import (
	"context"
	"iter"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

type Fn[T any] func(data T) bool

var _ pipe.Pipe = (*FilesPipe)(nil)

type FilesPipe struct {
	*common.BasePipe

	filters []Fn[*types.File]
}

func Files(filters ...Fn[*types.File]) *FilesPipe {
	return &FilesPipe{
		BasePipe: common.NewBase(),
		filters:  filters,
	}
}

func (p *FilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		for _, filter := range p.filters {
			file, ok := next()
			if !ok {
				break
			}

			if !filter(file) {
				continue
			}

			if !yield(file) {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
