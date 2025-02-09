package common

import (
	"context"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*BasePipe)(nil)

type BasePipe struct {
	next pipe.Pipe
	prev pipe.Pipe
}

func NewBase() *BasePipe {
	return &BasePipe{}
}

func (p *BasePipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	next := p.GetNext()
	if next == nil {
		return state, nil
	}

	done := make(chan struct{})

	var err error
	go func() {
		defer close(done)

		state, err = next.Run(ctx, state)
	}()

	select {
	case <-done:
		return state, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *BasePipe) GetNext() pipe.Pipe {
	return p.next
}

func (p *BasePipe) GetPrev() pipe.Pipe {
	return p.prev
}

func (p *BasePipe) SetNext(next pipe.Pipe) {
	p.next = next
}

func (p *BasePipe) SetPrev(prev pipe.Pipe) {
	p.prev = prev
}
