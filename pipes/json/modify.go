package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

type ModifyFn func(data map[string]any) error

var NopModifyFn = func(data map[string]any) error { return nil }

var _ pipe.Pipe = (*ModifyPipe)(nil)

type ModifyPipe struct {
	*common.BasePipe

	modifiers []ModifyFn
}

func Modify(modifiers ...ModifyFn) *ModifyPipe {
	return &ModifyPipe{
		BasePipe:  common.NewBase(),
		modifiers: modifiers,
	}
}

func (p *ModifyPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
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
		for _, modifier := range p.modifiers {
			if ok := func() bool {
				file, err, ok := next()
				if err != nil {
					return yield(nil, fmt.Errorf("json: modify: next: %w", err))
				}
				if !ok {
					return false
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
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("json: modify: unmarshal: %w", err))
						return
					}

					if err := modifier(jsonData); err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("json: modify: modifier: %w", err))
						return
					}

					if err := json.NewEncoder(pipeWriter).Encode(&jsonData); err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("json: modify: marshal: %w", err))
						return
					}
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
