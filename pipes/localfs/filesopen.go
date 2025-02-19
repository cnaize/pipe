package localfs

import (
	"context"
	"fmt"
	"io"
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

	var toClose []io.Closer
	defer func() {
		for _, closer := range toClose {
			closer.Close()
		}
	}()

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		for file := range files {
			if !yield(file) {
				break
			}
		}

		for _, name := range p.names {
			file, err := os.Open(name)
			if err != nil {
				syncErr.Join(fmt.Errorf("localfs: open files: open: %w", err))
				continue
			}
			toClose = append(toClose, file)

			stat, err := file.Stat()
			if err != nil {
				syncErr.Join(fmt.Errorf("localfs: open files: stat: %w", err))
				continue
			}

			if !yield(&types.File{
				Name: file.Name(),
				Data: file,
				Size: stat.Size(),
			}) {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
