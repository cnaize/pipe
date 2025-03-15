package tests

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/modify"
	"github.com/cnaize/pipe/pipes/state"
	"github.com/cnaize/pipe/types"
	"github.com/stretchr/testify/require"
)

type Data struct {
	Name    string `json:"name" yaml:"name"`
	Count   int    `json:"count" yaml:"count"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
}

func (suite *BaseTestSuite) TestModifyPipe() {
	inData0 := bytes.NewBufferString(`{
		"name": "json_0",
		"enabled": true
	}`)
	inData1 := bytes.NewBufferString(`{
		"name": "json_1",
		"count": 10
	}`)

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

	jsonModifyFn := func(data *Data) error {
		data.Count++
		data.Enabled = !data.Enabled

		return nil
	}

	jsonLine := pipes.Line(
		readLine,
		modify.Jsons(jsonModifyFn, jsonModifyFn),
		writeLine,
	)

	res, err := jsonLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	var i int
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(44), file.Size)
			require.EqualValues(suite.T(), `{"name":"json_0","count":1,"enabled":false}`+string('\n'), outData0.String())
		} else {
			require.Equal(suite.T(), int64(44), file.Size)
			require.EqualValues(suite.T(), `{"name":"json_1","count":11,"enabled":true}`+string('\n'), outData1.String())
		}

		i++
	}

	inData0.Reset()
	inData1.Reset()
	outData0.Reset()
	outData1.Reset()

	inData0.WriteString("name: yaml_0\nenabled: true")
	inData1.WriteString("name: yaml_1\ncount: 10")

	yamlModifyFn := func(data *Data) error {
		data.Count++
		data.Enabled = !data.Enabled

		return nil
	}

	yamlLine := pipes.Line(
		readLine,
		modify.Yamls(yamlModifyFn, yamlModifyFn),
		writeLine,
	)

	res, err = yamlLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	i = 0
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(37), file.Size)
			require.EqualValues(suite.T(), "name: yaml_0\ncount: 1\nenabled: false\n", outData0.String())
		} else {
			require.Equal(suite.T(), int64(37), file.Size)
			require.EqualValues(suite.T(), "name: yaml_1\ncount: 11\nenabled: true\n", outData1.String())
		}

		i++
	}

	inData0.Reset()
	inData1.Reset()
	outData0.Reset()
	outData1.Reset()

	inData0.WriteString("name: file_0 enabled: true")
	inData1.WriteString("name: file_1 count: 10")

	fileModifyFn := func(file *types.File) error {
		data, err := io.ReadAll(file.Data)
		if err != nil {
			return err
		}

		newData := strings.ReplaceAll(string(data), "enabled: true", "enabled: false")

		file.Data = strings.NewReader(newData)
		file.Size = int64(len(newData))

		return nil
	}

	fileLine := pipes.Line(
		readLine,
		modify.Files(fileModifyFn, fileModifyFn),
		writeLine,
	)

	res, err = fileLine.Run(context.Background(), nil)
	require.NoError(suite.T(), err)

	i = 0
	for file := range res.Files {
		require.NotEmpty(suite.T(), file.Size)

		if i == 0 {
			require.Equal(suite.T(), int64(27), file.Size)
			require.EqualValues(suite.T(), "name: file_0 enabled: false", outData0.String())
		} else {
			require.Equal(suite.T(), int64(22), file.Size)
			require.EqualValues(suite.T(), "name: file_1 count: 10", outData1.String())
		}

		i++
	}
}
