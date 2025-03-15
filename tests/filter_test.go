package tests

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/filter"
	"github.com/cnaize/pipe/pipes/state"
	"github.com/cnaize/pipe/types"
	"github.com/stretchr/testify/require"
)

func (suite *BaseTestSuite) TestFilteripe() {
	inData0 := bytes.NewBufferString(`{"name": "json_0","enabled": true}`)
	inData1 := bytes.NewBufferString(`{"name": "json_1","count": 10}`)

	readLine := pipes.Line(
		common.Timeout(time.Second),
		common.ReadFrom(inData0, inData1),
	)

	outData0 := bytes.NewBuffer(nil)
	outData1 := bytes.NewBuffer(nil)

	writeLine := pipes.Line(
		common.WriteTo(outData0, outData1),
		state.Consume(),
	)

	jsonFilterFn := func(data Data) bool {
		return data.Enabled
	}

	jsonLine := pipes.Line(
		readLine,
		filter.Jsons(jsonFilterFn, jsonFilterFn),
		writeLine,
	)

	res, err := jsonLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	var i int
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(43), file.Size)
			require.EqualValues(suite.T(), `{"name":"json_0","count":0,"enabled":true}`+string('\n'), outData0.String())
		} else {
			require.Equal(suite.T(), 0, file.Size)
			require.Empty(suite.T(), outData1.String())
		}

		i++
	}

	inData0.Reset()
	inData1.Reset()
	outData0.Reset()
	outData1.Reset()

	inData0.WriteString("name: yaml_0\nenabled: true")
	inData1.WriteString("name: yaml_1\ncount: 10")

	yamlFilterFn := func(data Data) bool {
		return data.Enabled
	}

	yamlLine := pipes.Line(
		readLine,
		filter.Yamls(yamlFilterFn, yamlFilterFn),
		writeLine,
	)

	res, err = yamlLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	i = 0
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(36), file.Size)
			require.EqualValues(suite.T(), "name: yaml_0\ncount: 0\nenabled: true\n", outData0.String())
		} else {
			require.Equal(suite.T(), 0, file.Size)
			require.Empty(suite.T(), outData1.String())
		}

		i++
	}

	inData0.Reset()
	inData1.Reset()
	outData0.Reset()
	outData1.Reset()

	inData0.WriteString("name: file_0 enabled: true")
	inData1.WriteString("name: file_1 count: 10")

	fileFilterFn := func(file *types.File) bool {
		data, err := io.ReadAll(file.Data)
		if err != nil {
			return false
		}

		if !bytes.Contains(data, []byte("enabled: true")) {
			return false
		}

		file.Data = bytes.NewReader(data)

		return true
	}

	fileLine := pipes.Line(
		readLine,
		filter.Files(fileFilterFn, fileFilterFn),
		writeLine,
	)

	res, err = fileLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	i = 0
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(26), file.Size)
			require.EqualValues(suite.T(), "name: file_0 enabled: true", outData0.String())
		} else {
			require.Equal(suite.T(), 0, file.Size)
			require.Empty(suite.T(), outData1.String())
		}

		i++
	}
}
