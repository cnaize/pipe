package pipe

import (
	"context"
	"fmt"
	"io"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*DiscardPipe)(nil)

type DiscardPipe struct {
	*basePipe
}

func Discard() *DiscardPipe {
	return &DiscardPipe{
		basePipe: newBasePipe(),
	}
}

func (p *DiscardPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	if _, err := io.Copy(io.Discard, in.Data); err != nil {
		return nil, fmt.Errorf("discard: copy: %w", err)
	}

	return p.basePipe.Send(ctx, in)
}
