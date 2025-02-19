package tests

import (
	"bytes"
	"context"
	"time"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/modify"
	"github.com/cnaize/pipe/pipes/state"
	"github.com/stretchr/testify/require"
)

func (suite *BaseTestSuite) TestJsonPipe() {
	jsonData0 := bytes.NewBufferString(`{
		"name": "json_0",
		"enabled": true
	}`)
	jsonData1 := bytes.NewBufferString(`{
		"name": "json_1",
		"count": 10
	}`)

	modifyFn := func(data map[string]any) error {
		if enabled, ok := data["enabled"]; ok {
			enabled := enabled.(bool)
			data["enabled"] = !enabled
		}

		return nil
	}

	outData0 := bytes.NewBuffer(nil)
	outData1 := bytes.NewBuffer(nil)

	pipeline, err := pipes.Line(
		common.Timeout(time.Second),
		common.ReadFrom(jsonData0, jsonData1),
		modify.Jsons(modifyFn, modifyFn),
		common.WriteTo(outData0, outData1),
		state.ConsumeFiles(),
	)
	require.NoError(suite.T(), err)

	res, err := pipeline.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	var i int
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(34), file.Size)
			require.EqualValues(suite.T(), `{"enabled":false,"name":"json_0"}`+string('\n'), outData0.String())
		} else {
			require.Equal(suite.T(), int64(29), file.Size)
			require.EqualValues(suite.T(), `{"count":10,"name":"json_1"}`+string('\n'), outData1.String())
		}

		i++
	}
}
