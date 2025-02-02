package pipe

import (
	"context"
	"time"

	"github.com/cnaize/pipes/types"
)

type TimeoutPipe struct {
	*basePipe

	timeout time.Duration
}

func Timeout(timeout time.Duration) *TimeoutPipe {
	return &TimeoutPipe{
		basePipe: newBasePipe(),
		timeout:  timeout,
	}
}

func (p *TimeoutPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.basePipe.Send(ctx, in)
}
