package filter

import (
	"context"
	"fmt"
	"io"
	"iter"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*YamlsPipe[any])(nil)

type YamlsPipe[T any] struct {
	*common.BasePipe

	filters []Fn[T]
}

func Yamls[T any](filters ...Fn[T]) *YamlsPipe[T] {
	return &YamlsPipe[T]{
		BasePipe: common.NewBase(),
		filters:  filters,
	}
}

func (p *YamlsPipe[T]) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		var wg sync.WaitGroup
		for _, filter := range p.filters {
			file, ok := next()
			if !ok {
				break
			}

			var data T
			if err := yaml.NewDecoder(file.Data).Decode(&data); err != nil {
				syncErr.Join(fmt.Errorf("filter: yamls: unmarshal: %w", err))
				return
			}

			if !filter(data) {
				continue
			}

			pipeReader, pipeWriter := io.Pipe()
			file.Data = pipeReader

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer pipeWriter.Close()

				if err := yaml.NewEncoder(pipeWriter).Encode(&data); err != nil {
					syncErr.Join(fmt.Errorf("filter: yamls: marshal: %w", err))
					return
				}
			}()

			if !yield(file) {
				break
			}
		}

		wg.Wait()
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
