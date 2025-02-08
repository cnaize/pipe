package localfs

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*CreateFilesPipe)(nil)

type CreateFilesPipe struct {
	*general.BasePipe

	names []string
}

func CreateFiles(names ...string) *CreateFilesPipe {
	return &CreateFilesPipe{
		BasePipe: general.NewBase(),
		names:    names,
	}
}

func (p *CreateFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		next, stop := iter.Pull2(files)
		defer stop()

		for _, name := range p.names {
			if ok := func() bool {
				file, err, ok := next()
				if err != nil {
					return yield(nil, fmt.Errorf("lfs: create files: next: %w", err))
				}
				if !ok {
					return false
				}

				newFile, err := os.Create(name)
				if err != nil {
					return yield(nil, fmt.Errorf("lfs: create files: create: %w", err))
				}
				defer newFile.Close()

				if _, err := io.Copy(newFile, file.Data); err != nil {
					return yield(nil, fmt.Errorf("lfs: create files: copy: %w", err))
				}

				stat, err := newFile.Stat()
				if err != nil {
					return yield(nil, fmt.Errorf("lfs: create files: stat: %w", err))
				}

				file.Name = newFile.Name()
				file.Perm = stat.Mode()
				file.Size = stat.Size()

				return yield(file, nil)
			}(); !ok {
				break
			}
		}
	}

	return p.BasePipe.Run(ctx, state)
}
