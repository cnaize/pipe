package modify

import (
	"context"
	"fmt"
	"iter"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

type FileFn func(file *types.File) error

var _ pipe.Pipe = (*FilesPipe)(nil)

type FilesPipe struct {
	*common.BasePipe

	modifiers []FileFn
}

func Files(modifiers ...FileFn) *FilesPipe {
	return &FilesPipe{
		BasePipe:  common.NewBase(),
		modifiers: modifiers,
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

		for _, modifier := range p.modifiers {
			file, ok := next()
			if !ok {
				break
			}

			if err := modifier(file); err != nil {
				syncErr.Join(fmt.Errorf("files: modify: modifier: %w", err))
				return
			}

			if !yield(file) {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
