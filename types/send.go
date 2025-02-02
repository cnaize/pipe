package types

import (
	"io"
	"os"
)

type SendIn struct {
	Read io.Reader
}

type SendOut struct {
	DirMakeAll *Dir
	FileOpen   *File
	FileCreate *File
	Sha256     *string
}

type Dir struct {
	Path string
	Perm os.FileMode
}

type File struct {
	Path string
}
