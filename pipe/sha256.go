package pipe

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"slices"

	"github.com/cnaize/pipes/types"
)

type Sha256Pipe struct {
	*basePipe

	hash []byte
}

func Sha256(hash []byte) *Sha256Pipe {
	return &Sha256Pipe{
		basePipe: newBasePipe(),
		hash:     hash,
	}
}

func (p *Sha256Pipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	data := in.Read

	hash := sha256.New()
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()

		io.Copy(io.MultiWriter(hash, writer), data)
	}()

	in.Read = reader

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	out.Sha256 = hash.Sum(nil)
	if p.hash != nil && !slices.Equal(out.Sha256, p.hash) {
		return nil, fmt.Errorf("sha256: check hash: invalid hash")
	}

	return out, nil
}
