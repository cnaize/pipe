package pipe

import (
	"context"
	"fmt"
	"os"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*FileOpenPipe)(nil)

type FileOpenPipe struct {
	*basePipe

	path string
}

func FileOpen(path string) *FileOpenPipe {
	return &FileOpenPipe{
		basePipe: newBasePipe(),
		path:     path,
	}
}

func (p *FileOpenPipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	file, err := os.Open(p.path)
	if err != nil {
		return nil, fmt.Errorf("file open: open: %w", err)
	}
	defer file.Close()

	in.Read = file

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	out.FileOpen = &types.File{
		Path: p.path,
	}

	return out, nil
}
