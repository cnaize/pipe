package pipe

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/cnaize/pipes/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	err := os.RemoveAll("testdata/tmp")
	require.NoError(suite.T(), err)
}

func (suite *BaseTestSuite) TestPipe() {
	pipeline, err := Line(
		Timeout(time.Second),
		OsRemoveAll("testdata/tmp"),
		DirMakeAll("testdata/tmp", os.ModePerm),
		FileOpen("testdata/test.txt"),
		Sha256(),
		FileCreate("testdata/tmp/test.txt"),
	)
	require.NoError(suite.T(), err)

	out, err := pipeline.Send(context.Background(), &types.SendIn{})
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), out)

	require.NotNil(suite.T(), out.Sha256)
	require.EqualValues(suite.T(),
		&types.File{Path: "testdata/test.txt"},
		out.FileOpen,
	)
	require.EqualValues(suite.T(),
		&types.File{Path: "testdata/tmp/test.txt"},
		out.FileCreate,
	)

	testFile, err := os.Open("testdata/test.txt")
	require.NoError(suite.T(), err)
	defer testFile.Close()
	testData, err := io.ReadAll(testFile)
	require.NoError(suite.T(), err)

	tmpFile, err := os.Open("testdata/tmp/test.txt")
	require.NoError(suite.T(), err)
	defer tmpFile.Close()
	tmpData, err := io.ReadAll(tmpFile)
	require.NoError(suite.T(), err)

	require.EqualValues(suite.T(), testData, tmpData)
}
