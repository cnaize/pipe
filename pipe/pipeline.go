package pipe

import (
	"context"
	"fmt"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*Pipeline)(nil)

type Pipeline struct {
	*basePipe

	line []pipes.Pipe
}

func Line(line ...pipes.Pipe) (*Pipeline, error) {
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

	return &Pipeline{
		line: line,
	}, nil
}

func (p *Pipeline) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	return p.line[0].Send(ctx, in)
}
