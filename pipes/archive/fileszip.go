package archive

import (
	"archive/zip"
	"context"
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

	var syncErr types.SyncError

	files := state.Files
	state.Files = func(yield func(*types.File) bool) {
		pipeReader, pipeWriter := io.Pipe()
		go func() {
			defer pipeWriter.Close()

			zipWriter := zip.NewWriter(pipeWriter)
			defer zipWriter.Close()

			for file := range files {
				err := func() error {
					fileData := file.Data
					file.Data = nil

					zipFile, err := zipWriter.Create(filepath.Base(file.Name))
					if err != nil {
						return fmt.Errorf("archive: zip files: create: %w", err)
					}

					if _, err := io.Copy(zipFile, fileData); err != nil {
						return fmt.Errorf("archive: zip files: copy: %w", err)
					}

					return nil
				}()
				syncErr.Join(err)
			}
		}()

		yield(&types.File{Data: pipeReader})
	}

	state, err := p.BasePipe.Run(ctx, state)
	return state, syncErr.Join(err)
}
