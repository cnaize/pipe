package pipe

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/cnaize/pipes/types"
)

type Sha256Pipe struct {
	*basePipe

	expected *string
}

func Sha256(expected *string) *Sha256Pipe {
	return &Sha256Pipe{
		basePipe: newBasePipe(),
		expected: expected,
	}
}

func (p *Sha256Pipe) Send(ctx context.Context, in *types.SendIn) (*types.SendOut, error) {
	data := in.Read

	hasher := sha256.New()
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()

		io.Copy(io.MultiWriter(hasher, writer), data)
	}()

	in.Read = reader

	out, err := p.basePipe.Send(ctx, in)
	if err != nil {
		return nil, err
	}

	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	out.Sha256 = &hash

	if p.expected != nil && hash != *p.expected {
		return nil, fmt.Errorf("sha256: check hash: %w", types.ErrInvalidHash)
	}

	return out, nil
}
