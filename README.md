# Pipe
## Repo philosophy — everything is a pipe

### WARNING: very early stage, contibutors are welcome

## Example:
```go
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
	// craeate a directories pipeline
	dirsLine, _ := pipes.Line(
		// remove temporary directory
		localfs.RemoveAll("testdata/tmp"),
		// create temporary directory
		localfs.MakeDirAll("testdata/tmp", os.ModePerm),
	)

	// craeate a files pipeline
	filesLine, _ := pipes.Line(
		// set execution timeout
		common.Timeout(time.Second),
		// open two example files
		localfs.OpenFiles("testdata/test_0.txt", "testdata/test_1.txt"),
		// calculate and compare hash for each file
		hash.SumSha256("kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=", "CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y="),
		// zip the files
		archive.ZipFiles(),
		// calculate hash for the zip archive
		hash.SumSha256(""),
		// create a new file
		localfs.CreateFiles("testdata/tmp/test.zip"),
		// flow the files through the pipes and keep metadata
		state.Consume(),
	)

	// create a composite pipeline using the pipelines above
	pipeline, _ := pipes.Line(
		// add the directories pipeline
		dirsLine,
		// add the files pipeline
		filesLine,
	)

	// run the composite pipeline
	res, _ := pipeline.Run(context.Background(), nil)

	// iterate over result files and print metadata
	for file := range res.Files {
		fmt.Printf("Result file:\n\tName: %s\n\tSize: %d\n\tHash: %s\n", file.Name, file.Size, file.Hash)
	}
}
```
### Output:
```
Result file:
	Name: testdata/tmp/test.zip
	Size: 1047
	Hash: Yg3OOaBD-miLs7lDIBVAeZMZIXYfy2N25f8-b-1kWOc=
```
