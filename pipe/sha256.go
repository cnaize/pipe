package pipe

import (
	"context"
	"crypto/sha256"
	"io"

	"github.com/cnaize/pipes/types"
)

type Sha256Pipe struct {
	*basePipe
}

func Sha256() *Sha256Pipe {
	return &Sha256Pipe{
		basePipe: newBasePipe(),
	}
}

func (p *Sha256Pipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	data := in.Data

	hash := sha256.New()
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()

		io.Copy(io.MultiWriter(hash, writer), data)
	}()

	in.Data = reader

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	out.Sha256 = hash.Sum(nil)

	return out, nil
}
