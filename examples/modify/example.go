package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/cnaize/pipe/pipes"
	"github.com/cnaize/pipe/pipes/common"
	"github.com/cnaize/pipe/pipes/modify"
	"github.com/cnaize/pipe/pipes/state"
)

func main() {
	// create two example jsons
	inData0 := bytes.NewBufferString(`{
		"name": "json_0",
		"enabled": true
	}`)
	inData1 := bytes.NewBufferString(`{
		"name": "json_1",
		"count": 10
	}`)

	// create read pipeline
	readLine := pipes.Line(
		// set execution timeout
		common.Timeout(time.Second),
		// pass the example inputs
		common.ReadFrom(inData0, inData1),
	)

	// create output buffers
	outData0 := bytes.NewBuffer(nil)
	outData1 := bytes.NewBuffer(nil)

	// create write pipeline
	writeLine := pipes.Line(
		// pass the output buffers
		common.WriteTo(outData0, outData1),
		// flow the jsons through the pipes and keep metadata
		state.Consume(),
	)

	// create json modify function
	jsonModifyFn := func(data map[string]any) error {
		if enabled, ok := data["enabled"]; ok {
			enabled := enabled.(bool)
			data["enabled"] = !enabled
		}

		return nil
	}

	// craeate json pipeline
	jsonLine := pipes.Line(
		// pass the read pipeline
		readLine,
		// pass the json modify function
		modify.Jsons(jsonModifyFn, jsonModifyFn),
		// pass the write pipeline
		writeLine,
	)

	// run the json pipeline
	_, _ = jsonLine.Run(context.Background(), nil)

	// print the output buffers data
	fmt.Printf("=== Json ===\n\n")
	fmt.Printf("--> Result data 0:\n%s\n", outData0.String())
	fmt.Printf("--> Redult data 1:\n%s\n", outData1.String())

	// reset the buffers
	inData0.Reset()
	inData1.Reset()
	outData0.Reset()
	outData1.Reset()

	// create two example yamls
	inData0.WriteString("name: yaml_0\nenabled: true")
	inData1.WriteString("name: yaml_1\ncount: 10")

	// create yaml modify function
	yamlModifyFn := func(data map[any]any) error {
		if enabled, ok := data["enabled"]; ok {
			enabled := enabled.(bool)
			data["enabled"] = !enabled
		}

		return nil
	}

	// craeate yaml pipeline
	yamlLine := pipes.Line(
		// pass the read pipeline
		readLine,
		// pass the json modify function
		modify.Yamls(yamlModifyFn, yamlModifyFn),
		// pass the write pipeline
		writeLine,
	)

	// run the yaml pipeline
	_, _ = yamlLine.Run(context.Background(), nil)

	// print the output buffers data
	fmt.Printf("=== Yaml ===\n\n")
	fmt.Printf("--> Result data 0:\n%s\n", outData0.String())
	fmt.Printf("--> Redult data 1:\n%s\n", outData1.String())
}
