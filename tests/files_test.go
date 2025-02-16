package tests

import (
	"context"
	"os"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/archive"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/hash"
	"github.com/cnaize/pipe/pipes/localfs"
	"github.com/cnaize/pipe/pipes/state"
	"github.com/cnaize/pipe/types"
)

func (suite *BaseTestSuite) TestFilesPipe() {
	dirsLine, err := pipes.Line(
		localfs.RemoveAll("../testdata/tmp"),
		localfs.MakeDirAll("../testdata/tmp", os.ModePerm),
	)
	require.NoError(suite.T(), err)

	filesLine, err := pipes.Line(
		common.Timeout(time.Second),
		localfs.OpenFiles("../testdata/test_0.txt", "../testdata/test_1.txt"),
		hash.SumSha256("kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=", "CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y="),
		archive.ZipFiles(),
		hash.SumSha256(""),
		localfs.CreateFiles("../testdata/tmp/test.zip"),
		state.ConsumeFiles(),
	)
	require.NoError(suite.T(), err)

	pipeline, err := pipes.Line(
		dirsLine,
		filesLine,
	)
	require.NoError(suite.T(), err)

	res, err := pipeline.Run(context.Background(), nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	for file := range res.Files {
		require.EqualValues(suite.T(),
			&types.File{
				Name: "../testdata/tmp/test.zip",
				Size: 1047,
				Hash: "Yg3OOaBD-miLs7lDIBVAeZMZIXYfy2N25f8-b-1kWOc=",
			},
			file,
		)
	}
}
