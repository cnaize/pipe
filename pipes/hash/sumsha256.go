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
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*SumSha256Pipe)(nil)

type SumSha256Pipe struct {
	*common.BasePipe

	expected []string
}

func SumSha256(expected ...string) *SumSha256Pipe {
	return &SumSha256Pipe{
		BasePipe: common.NewBase(),
		expected: expected,
	}
}

func (p *SumSha256Pipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
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
					return yield(nil, fmt.Errorf("hash: sum sha256: next: %w", err))
				}
				if !ok {
					return false
				}

				hasher := sha256.New()
				pipeReader, pipeWriter := io.Pipe()

				fileData := file.Data
				file.Data = pipeReader

				go func() {
					defer pipeWriter.Close()

					if _, err := io.Copy(io.MultiWriter(hasher, pipeWriter), fileData); err != nil {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("hash: sum sha256: copy: %w", err))
						return
					}

					file.Hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
					if expected != "" && file.Hash != expected {
						mu.Lock()
						defer mu.Unlock()

						errs = errors.Join(errs, fmt.Errorf("hash: sum sha256: check: %w", types.ErrInvalidHash))
						return
					}
				}()

				return yield(file, nil)
			}(); !ok {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
