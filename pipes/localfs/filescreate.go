package localfs

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*CreateFilesPipe)(nil)

type CreateFilesPipe struct {
	*common.BasePipe

	names []string
}

func CreateFiles(names ...string) *CreateFilesPipe {
	return &CreateFilesPipe{
		BasePipe: common.NewBase(),
		names:    names,
	}
}

func (p *CreateFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		var wg sync.WaitGroup
		for _, name := range p.names {
			ok, err := func() (bool, error) {
				file, ok := next()
				if !ok {
					return false, nil
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					newFile, err := os.Create(name)
					if err != nil {
						syncErr.Join(fmt.Errorf("localfs: create files: create: %w", err))
						return
					}
					defer newFile.Close()

					if _, err := io.Copy(newFile, fileData); err != nil {
						syncErr.Join(fmt.Errorf("localfs: create files: copy: %w", err))
						return
					}

					stat, err := newFile.Stat()
					if err != nil {
						syncErr.Join(fmt.Errorf("localfs: create files: stat: %w", err))
						return
					}

					file.Name = newFile.Name()
					file.Size = stat.Size()
				}()

				return yield(file), nil
			}()
			syncErr.Join(err)
			if !ok {
				break
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
