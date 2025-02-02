package types

import "errors"

var (
	ErrInvalidHash   = errors.New("invalid hash")
	ErrEmptyPipeline = errors.New("empty pipeline")
)
