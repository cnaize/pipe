package pipe

import (
	"context"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*NopPipe)(nil)

type NopPipe struct {
	*basePipe
}

func Nop() *NopPipe {
	return &NopPipe{
		basePipe: newBasePipe(),
	}
}

func (p *NopPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	return p.basePipe.Send(ctx, in)
}
