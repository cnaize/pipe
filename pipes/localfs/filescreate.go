package localfs

import (
	"context"
	"errors"
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

	var mu sync.Mutex
	var errs error

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		next, stop := iter.Pull2(files)
		defer stop()

		var wg sync.WaitGroup
		for _, name := range p.names {
			if ok := func() bool {
				file, err, ok := next()
				if err != nil {
					return yield(nil, fmt.Errorf("localfs: create files: next: %w", err))
				}
				if !ok {
					return false
				}

				fileData := file.Data
				file.Data = nil

				wg.Add(1)
				go func() {
					defer wg.Done()

					newFile, err := os.Create(name)
					if err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("localfs: create files: create: %w", err))
						return
					}
					defer newFile.Close()

					if _, err := io.Copy(newFile, fileData); err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("localfs: create files: copy: %w", err))
						return
					}

					stat, err := newFile.Stat()
					if err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("localfs: create files: stat: %w", err))
						return
					}

					file.Name = newFile.Name()
					file.Perm = stat.Mode()
					file.Size = stat.Size()
				}()

				return yield(file, nil)
			}(); !ok {
				break
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
