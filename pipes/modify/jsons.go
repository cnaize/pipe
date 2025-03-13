package modify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

type JsonFn func(data map[string]any) error

var _ pipe.Pipe = (*JsonsPipe)(nil)

type JsonsPipe struct {
	*common.BasePipe

	modifiers []JsonFn
}

func Jsons(modifiers ...JsonFn) *JsonsPipe {
	return &JsonsPipe{
		BasePipe:  common.NewBase(),
		modifiers: modifiers,
	}
}

func (p *JsonsPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
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

				var jsonData map[string]any
				if err := json.NewDecoder(fileData).Decode(&jsonData); err != nil {
					syncErr.Join(fmt.Errorf("json: modify: unmarshal: %w", err))
					return
				}

				if err := modifier(jsonData); err != nil {
					syncErr.Join(fmt.Errorf("json: modify: modifier: %w", err))
					return
				}

				if err := json.NewEncoder(pipeWriter).Encode(&jsonData); err != nil {
					syncErr.Join(fmt.Errorf("json: modify: marshal: %w", err))
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
