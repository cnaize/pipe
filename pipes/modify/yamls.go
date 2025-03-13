package modify

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

type YamlFn func(data map[any]any) error

var _ pipe.Pipe = (*YamlsPipe)(nil)

type YamlsPipe struct {
	*common.BasePipe

	modifiers []YamlFn
}

func Yamls(modifiers ...YamlFn) *YamlsPipe {
	return &YamlsPipe{
		BasePipe:  common.NewBase(),
		modifiers: modifiers,
	}
}

func (p *YamlsPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		var wg sync.WaitGroup
		for _, modifier := range p.modifiers {
			file, ok := next()
			if !ok {
				break
			}

			pipeReader, pipeWriter := io.Pipe()

			fileData := file.Data
			file.Data = pipeReader

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer pipeWriter.Close()

				var yamlData map[any]any
				if err := yaml.NewDecoder(fileData).Decode(&yamlData); err != nil {
					syncErr.Join(fmt.Errorf("yaml: modify: unmarshal: %w", err))
					return
				}

				if err := modifier(yamlData); err != nil {
					syncErr.Join(fmt.Errorf("yaml: modify: modifier: %w", err))
					return
				}

				if err := yaml.NewEncoder(pipeWriter).Encode(&yamlData); err != nil {
					syncErr.Join(fmt.Errorf("yaml: modify: marshal: %w", err))
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
