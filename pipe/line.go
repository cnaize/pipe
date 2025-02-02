package pipe

import (
	"context"
	"fmt"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*LinePipe)(nil)

type LinePipe struct {
	*basePipe

	line []pipes.Pipe
}

func Line(line ...pipes.Pipe) (*LinePipe, error) {
	if len(line) < 1 {
		return nil, fmt.Errorf("empty line")
	}

	var prev pipes.Pipe
	for _, p := range line {
		p.SetPrev(prev)
		if prev != nil {
			prev.SetNext(p)
		}

		prev = p
	}

	return &LinePipe{
		basePipe: newBasePipe(),
		line:     line,
	}, nil
}

func (p *LinePipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	lineOut, err := p.line[0].Send(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("line: send: %w", err)
	}

	pipeOut, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	return types.MergeSendOut(lineOut, pipeOut), nil
}
