package hash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/general"
	pstate "github.com/cnaize/pipe/pipes/state"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*Sha256SumPipe)(nil)

type Sha256SumPipe struct {
	*general.BasePipe

	expected []string
}

func Sha256Sum(expected ...string) *Sha256SumPipe {
	return &Sha256SumPipe{
		BasePipe: general.NewBase(),
		expected: expected,
	}
}

func (p *Sha256SumPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var mu sync.Mutex
	var errs error

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		next, stop := iter.Pull2(files)
		defer stop()

		for _, expected := range p.expected {
			if ok := func() bool {
				file, err, ok := next()
				if err != nil {
					return yield(nil, fmt.Errorf("hash: sha256 sum: next: %w", err))
				}
				if !ok {
					return false
				}

				data := file.Data
				hasher := sha256.New()
				pipeReader, pipeWriter := io.Pipe()
				go func() {
					defer pipeWriter.Close()

					if _, err := io.Copy(io.MultiWriter(hasher, pipeWriter), data); err != nil {
						errs = errors.Join(errs, fmt.Errorf("hash: sha256 sum: copy: %w", err))
					}

					hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
					if expected != "" && hash != expected {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("hash: sha256 sum: check: %w", types.ErrInvalidHash))
					}
					file.Hash = &hash
				}()

				file.Data = pipeReader

				return yield(file, nil)
			}(); !ok {
				break
			}
		}
	}

	if p.GetNext() == nil {
		var err error
		if state, err = pstate.DiscardFiles().Run(ctx, state); err != nil {
			return nil, errors.Join(errs, fmt.Errorf("hash: sha256 sum: discard files: %w", err))
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
