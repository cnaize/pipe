# Pipe repo philosophy: everything is a pipe

### WARNING: very early stage, contributions welcome

## Example:
```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cnaize/pipe/pipes"
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
		// calculate hash for the first file and calculate and compare hash for the second file
		hash.SumSha256("", "CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y="),
		// create two new files
		localfs.CreateFiles("testdata/tmp/test_0.txt", "testdata/tmp/test_1.txt"),
		// run the files through pipes and keep the files metadata
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
		fmt.Printf("Result file:\n\tName: %s\n\tSize: %d\n\tHash: %s\n", file.Name, file.Size, file.Hash)
	}
}
```
### Output:
```
Result file:
        Name: testdata/tmp/test_0.txt
        Size: 502
        Hash: kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=
Result file:
        Name: testdata/tmp/test_1.txt
        Size: 946
        Hash: CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y=
```