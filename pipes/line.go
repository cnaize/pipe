package pipes

import (
	"context"
	"fmt"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*LinePipe)(nil)

type LinePipe struct {
	*common.BasePipe

	line []pipe.Pipe
}

func Line(line ...pipe.Pipe) *LinePipe {
	var prev pipe.Pipe
	for _, p := range line {
		p.SetPrev(prev)
		if prev != nil {
			prev.SetNext(p)
		}

		prev = p
	}

	return &LinePipe{
		BasePipe: common.NewBase(),
		line:     line,
	}
}

func (p *LinePipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if len(p.line) > 0 {
		var err error
		
		state, err = p.line[0].Run(ctx, state)
		if err != nil {
			return nil, fmt.Errorf("line: run: %w", err)
		}
	}

	return p.BasePipe.Run(ctx, state)
}
