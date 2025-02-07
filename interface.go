package pipe

import (
	"context"

	"github.com/cnaize/pipe/types"
)

type Pipe interface {
	Run(ctx context.Context, state *types.State) (*types.State, error)

	GetNext() Pipe
	GetPrev() Pipe

	SetNext(next Pipe)
	SetPrev(prev Pipe)
}
