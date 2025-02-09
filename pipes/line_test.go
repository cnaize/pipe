package pipe

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/hash"
	"github.com/cnaize/pipe/pipes/localfs"
	"github.com/cnaize/pipe/pipes/state"
	"github.com/cnaize/pipe/types"
)

type BaseTestSuite struct {
	suite.Suite
}

func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}

func (suite *BaseTestSuite) SetupSuite() {

}

func (suite *BaseTestSuite) TearDownSuite() {
	err := os.RemoveAll("../testdata/tmp")
	require.NoError(suite.T(), err)
}

func (suite *BaseTestSuite) TestPipe() {
	dirsLine, err := Line(
		localfs.RemoveAll("../testdata/tmp"),
		localfs.MakeDirAll("../testdata/tmp", os.ModePerm),
	)
	require.NoError(suite.T(), err)

	filesLine, err := Line(
		common.Timeout(time.Second),
		localfs.OpenFiles("../testdata/test_0.txt", "../testdata/test_1.txt"),
		hash.SumSha256("kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=", ""),
		localfs.CreateFiles("../testdata/tmp/test_0.txt", "../testdata/tmp/test_1.txt"),
		state.ConsumeFiles(),
	)
	require.NoError(suite.T(), err)

	pipeline, err := Line(
		dirsLine,
		filesLine,
	)
	require.NoError(suite.T(), err)

	res, err := pipeline.Run(context.Background(), nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)

	var i int
	for file, err := range res.Files {
		require.NoError(suite.T(), err)

		if i == 0 {
			require.EqualValues(suite.T(),
				&types.File{
					Name: fmt.Sprintf("../testdata/tmp/test_%d.txt", i),
					Perm: 0644,
					Size: 502,
					Hash: "kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c=",
				},
				file,
			)
		} else {
			require.EqualValues(suite.T(),
				&types.File{
					Name: fmt.Sprintf("../testdata/tmp/test_%d.txt", i),
					Perm: 0644,
					Size: 946,
					Hash: "CeE_WA_xKsx2Dj_sRvowaCeDfQOPviSpyjaZdxuCT4Y=",
				},
				file,
			)
		}

		testFile, err := os.Open(fmt.Sprintf("../testdata/test_%d.txt", i))
		require.NoError(suite.T(), err)
		defer testFile.Close()
		testData, err := io.ReadAll(testFile)
		require.NoError(suite.T(), err)

		tmpFile, err := os.Open(fmt.Sprintf("../testdata/tmp/test_%d.txt", i))
		require.NoError(suite.T(), err)
		defer tmpFile.Close()
		tmpData, err := io.ReadAll(tmpFile)
		require.NoError(suite.T(), err)

		require.EqualValues(suite.T(), testData, tmpData)

		i++
	}
}
