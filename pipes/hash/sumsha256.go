package hash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"iter"

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

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		for _, expected := range p.expected {
			ok, err := func() (bool, error) {
				file, ok := next()
				if !ok {
					return false, nil
				}

				hasher := sha256.New()
				pipeReader, pipeWriter := io.Pipe()

				fileData := file.Data
				file.Data = pipeReader

				go func() {
					defer pipeWriter.Close()

					if _, err := io.Copy(io.MultiWriter(hasher, pipeWriter), fileData); err != nil {
						syncErr.Join(fmt.Errorf("hash: sum sha256: copy: %w", err))
						return
					}

					file.Hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
					if expected != "" && file.Hash != expected {
						syncErr.Join(fmt.Errorf("hash: sum sha256: check: %w", types.ErrInvalidHash))
						return
					}
				}()

				return yield(file), nil
			}()
			syncErr.Join(err)
			if !ok {
				break
			}
		}
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
