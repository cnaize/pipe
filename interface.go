package pipes

import (
	"context"

	"github.com/cnaize/pipes/types"
)

type Pipe interface {
	Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error)

	GetNext() Pipe
	GetPrev() Pipe

	SetNext(next Pipe)
	SetPrev(prev Pipe)
}
