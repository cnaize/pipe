package pipe

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cnaize/pipes"
	"github.com/cnaize/pipes/types"
)

var _ pipes.Pipe = (*FileCreatePipe)(nil)

type FileCreatePipe struct {
	*basePipe

	path string
}

func FileCreate(path string) *FileCreatePipe {
	return &FileCreatePipe{
		basePipe: newBasePipe(),
		path:     path,
	}
}

func (p *FileCreatePipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	file, err := os.Create(p.path)
	if err != nil {
		return nil, fmt.Errorf("file create: create: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, in.Data); err != nil {
		return nil, fmt.Errorf("file create: copy: %w", err)
	}

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	out.FileCreate = &types.File{
		Path: p.path,
	}

	return out, nil
}
