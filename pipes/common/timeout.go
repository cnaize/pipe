package common

import (
	"context"
	"time"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*TimeoutPipe)(nil)

type TimeoutPipe struct {
	*BasePipe

	timeout time.Duration
}

func Timeout(timeout time.Duration) *TimeoutPipe {
	return &TimeoutPipe{
		BasePipe: NewBase(),
		timeout:  timeout,
	}
}

func (p *TimeoutPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.BasePipe.Run(ctx, state)
}
