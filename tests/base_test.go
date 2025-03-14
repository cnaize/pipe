package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Data struct {
	Name    string `json:"name" yaml:"name"`
	Count   int    `json:"count" yaml:"count"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
}

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
