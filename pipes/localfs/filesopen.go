package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*OpenFilesPipe)(nil)

type OpenFilesPipe struct {
	*common.BasePipe

	names []string
}

func OpenFiles(names ...string) *OpenFilesPipe {
	return &OpenFilesPipe{
		BasePipe: common.NewBase(),
		names:    names,
	}
}

func (p *OpenFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		for file := range files {
			if !yield(file) {
				break
			}
		}

		for _, name := range p.names {
			ok, err := func() (bool, error) {
				file, err := os.Open(name)
				if err != nil {
					return false, fmt.Errorf("localfs: open files: open: %w", err)
				}
				defer file.Close()

				stat, err := file.Stat()
				if err != nil {
					return false, fmt.Errorf("localfs: open files: stat: %w", err)
				}

				return yield(&types.File{
					Name: file.Name(),
					Data: file,
					Size: stat.Size(),
				}), nil
			}()
			syncErr.Join(err)
			if !ok {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
