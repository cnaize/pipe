package types

import (
	"io"
	"iter"
	"os"
)

type File struct {
	Name string
	Perm os.FileMode
	Data io.Reader
	Size int64
	Hash string
}

type State struct {
	Files iter.Seq2[*File, error]
}

func NewState() *State {
	return &State{
		Files: func(yield func(*File, error) bool) {},
	}
}
