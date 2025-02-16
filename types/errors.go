package types

import (
	"errors"
	"sync"
)

var (
	ErrInvalidHash   = errors.New("invalid hash")
	ErrEmptyPipeline = errors.New("empty pipeline")
)

type SyncError struct {
	sync.Mutex
	error
}

func (e *SyncError) Join(err error) error {
	e.Lock()
	defer e.Unlock()

	e.error = errors.Join(e.error, err)
	return e.error
}
