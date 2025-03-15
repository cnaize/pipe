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

var _ pipe.Pipe = (*JsonsPipe[any])(nil)

type JsonsPipe[T any] struct {
	*common.BasePipe

	modifiers []Fn[T]
}

func Jsons[T any](modifiers ...Fn[T]) *JsonsPipe[T] {
	return &JsonsPipe[T]{
		BasePipe:  common.NewBase(),
		modifiers: modifiers,
	}
}

func (p *JsonsPipe[T]) Run(ctx context.Context, state *types.State) (*types.State, error) {
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

				var data T
				if err := json.NewDecoder(fileData).Decode(&data); err != nil {
					syncErr.Join(fmt.Errorf("json: modify: unmarshal: %w", err))
					return
				}

				if err := modifier(data); err != nil {
					syncErr.Join(fmt.Errorf("json: modify: modifier: %w", err))
					return
				}

				if err := json.NewEncoder(pipeWriter).Encode(&data); err != nil {
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
