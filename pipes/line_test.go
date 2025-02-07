package pipe

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cnaize/pipe/pipes/general"
	"github.com/cnaize/pipe/pipes/hash"
	lfs "github.com/cnaize/pipe/pipes/localfs"
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
	dirLine, err := Line(
		lfs.OSRemoveAll("../testdata/tmp"),
		lfs.MakeDirAll("../testdata/tmp", os.ModePerm),
	)
	require.NoError(suite.T(), err)

	fileLine, err := Line(
		general.Timeout(time.Second),
		lfs.OpenFiles("../testdata/test.txt"),
		hash.Sha256Sum("kEvuni09HxM1ox-0nIj7_Ug1Adw0oIU62ukuh49oi5c="),
		lfs.CreateFiles("../testdata/tmp/test.txt"),
	)
	require.NoError(suite.T(), err)

	pipeline, err := Line(
		dirLine,
		fileLine,
	)
	require.NoError(suite.T(), err)

	state, err := pipeline.Run(context.Background(), nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), state)

	// require.NotNil(suite.T(), out.Sha256)
	// require.EqualValues(suite.T(), hash, *out.Sha256)
	// require.EqualValues(suite.T(),
	// 	&types.File{Path: "../testdata/test.txt"},
	// 	out.FileOpen,
	// )
	// require.EqualValues(suite.T(),
	// 	&types.File{Path: "../testdata/tmp/test.txt"},
	// 	out.FileCreate,
	// )

	testFile, err := os.Open("../testdata/test.txt")
	require.NoError(suite.T(), err)
	defer testFile.Close()
	testData, err := io.ReadAll(testFile)
	require.NoError(suite.T(), err)

	tmpFile, err := os.Open("../testdata/tmp/test.txt")
	require.NoError(suite.T(), err)
	defer tmpFile.Close()
	tmpData, err := io.ReadAll(tmpFile)
	require.NoError(suite.T(), err)

	require.EqualValues(suite.T(), testData, tmpData)
}
