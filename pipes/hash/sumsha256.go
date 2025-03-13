package hash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		next, stop := iter.Pull(files)
		defer stop()

		var wg sync.WaitGroup
		for _, expected := range p.expected {
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

				hasher := sha256.New()
				_, err := io.Copy(io.MultiWriter(hasher, pipeWriter), fileData)
				if err != nil {
					syncErr.Join(fmt.Errorf("hash: sum sha256: copy: %w", err))
					return
				}

				file.Hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

				if expected != "" && file.Hash != expected {
					syncErr.Join(fmt.Errorf("hash: sum sha256: check: %w", types.ErrInvalidHash))
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
