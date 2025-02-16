package types

import (
	"io"
	"iter"
)

type File struct {
	Name string
	Data io.Reader
	Size int64
	Hash string
}

type State struct {
	Files iter.Seq[*File]
}

func NewState() *State {
	return &State{
		Files: func(yield func(*File) bool) {},
	}
}
