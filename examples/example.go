package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/archive"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/hash"
	"github.com/cnaize/pipe/pipes/localfs"
	"github.com/cnaize/pipe/pipes/state"
)

func main() {
	// craeate directories pipeline
	dirsLine, _ := pipes.Line(
		// remove temporary directory
		localfs.RemoveAll("testdata/tmp"),
		// create temporary directory
		localfs.MakeDirAll("testdata/tmp", os.ModePerm),
	)

	// craeate files pipeline
	filesLine, _ := pipes.Line(
		// set execution timeout
		common.Timeout(time.Second),
		// open two example files
		localfs.OpenFiles("testdata/test_0.txt", "testdata/test_1.txt"),
		// calculate and compare hash for each file
		hash.SumSha256("kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=", "CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y="),
		// zip files
		archive.ZipFiles(),
		// create new file
		localfs.CreateFiles("testdata/tmp/test.zip"),
		// flow the files through the pipes and keep metadata
		state.ConsumeFiles(),
	)

	// create composite pipeline using the pipelines above
	pipeline, _ := pipes.Line(
		// add directories pipeline
		dirsLine,
		// add files pipeline
		filesLine,
	)

	// run the composite pipeline
	res, _ := pipeline.Run(context.Background(), nil)

	// iterate over result files
	for file, _ := range res.Files {
		fmt.Printf("Result file:\n\tName: %s\n\tSize: %d\n", file.Name, file.Size)
	}
}
