package localfs

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*OpenFilesPipe)(nil)

type OpenFilesPipe struct {
	*general.BasePipe

	names []string
}

func OpenFiles(names ...string) *OpenFilesPipe {
	return &OpenFilesPipe{
		BasePipe: general.NewBase(),
		names:    names,
	}
}

func (p *OpenFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		// get files
		for file, err := range files {
			if !yield(file, err) {
				break
			}
		}

		// add files
		for _, name := range p.names {
			if ok := func() bool {
				file, err := os.Open(name)
				if err != nil {
					return yield(nil, fmt.Errorf("lfs: open files: open: %w", err))
				}
				defer file.Close()

				stat, err := file.Stat()
				if err != nil {
					return yield(nil, fmt.Errorf("lfs: open files: stat: %w", err))
				}

				return yield(&types.File{
					Name: file.Name(),
					Perm: stat.Mode(),
					Data: file,
					Size: stat.Size(),
				}, nil)
			}(); !ok {
				break
			}
		}
	}

	return p.BasePipe.Run(ctx, state)
}
