package pipe

import (
	"context"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*basePipe)(nil)

type basePipe struct {
	prev pipes.Pipe
	next pipes.Pipe
}

func newBasePipe() *basePipe {
	return &basePipe{}
}

func (p *basePipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	next := p.GetNext()
	if next == nil {
		return &types.SendOut{}, nil
	}

	done := make(chan struct{})

	var out *types.SendOut
	var err error
	go func() {
		defer close(done)

		out, err = next.Send(ctx, in)
	}()

	select {
	case <-done:
		return out, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *basePipe) GetNext() pipes.Pipe {
	return p.next
}

func (p *basePipe) GetPrev() pipes.Pipe {
	return p.prev
}

func (p *basePipe) SetNext(next pipes.Pipe) {
	p.next = next
}

func (p *basePipe) SetPrev(prev pipes.Pipe) {
	p.prev = prev
}
