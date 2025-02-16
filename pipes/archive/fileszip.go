package archive

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/cnaize/pipe"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/types"
)

var _ pipe.Pipe = (*ZipFilesPipe)(nil)

type ZipFilesPipe struct {
	*common.BasePipe
}

func ZipFiles() *ZipFilesPipe {
	return &ZipFilesPipe{
		BasePipe: common.NewBase(),
	}
}

func (p *ZipFilesPipe) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state == nil {
		state = types.NewState()
	}

	var errs error

	files := state.Files
	state.Files = func(yield func(*types.File, error) bool) {
		pipeReader, pipeWriter := io.Pipe()
		go func() {
			defer pipeWriter.Close()

			zipWriter := zip.NewWriter(pipeWriter)
			defer zipWriter.Close()

			for file, err := range files {
				if err != nil {
					errs = errors.Join(errs, fmt.Errorf("archive: zip files: next: %w", err))
					return
				}

				fileData := file.Data
				file.Data = nil

				zipFile, err := zipWriter.Create(filepath.Base(file.Name))
				if err != nil {
					errs = errors.Join(errs, fmt.Errorf("archive: zip files: create: %w", err))
					return
				}

				if _, err := io.Copy(zipFile, fileData); err != nil {
					errs = errors.Join(errs, fmt.Errorf("archive: zip files: copy: %w", err))
					return
				}
			}
		}()

		yield(&types.File{
			Data: pipeReader,
		}, nil)
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, errors.Join(errs, err)
}
