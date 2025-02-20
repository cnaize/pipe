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
	jsonData0 := bytes.NewBufferString(`{
		"name": "json0",
		"enabled": true
	}`)
	jsonData1 := bytes.NewBufferString(`{
		"name": "json1",
		"count": 10
	}`)

	// create json modify function
	modifyFn := func(data map[string]any) error {
		if enabled, ok := data["enabled"]; ok {
			enabled := enabled.(bool)
			data["enabled"] = !enabled
		}

		return nil
	}

	// create output buffers
	outData0 := bytes.NewBuffer(nil)
	outData1 := bytes.NewBuffer(nil)

	// craeate a pipeline
	pipeline, _ := pipes.Line(
		// set execution timeout
		common.Timeout(time.Second),
		// pass the example jsons
		common.ReadFrom(jsonData0, jsonData1),
		// pass json modify functions
		modify.Jsons(modifyFn, modifyFn),
		// pass the output buffers
		common.WriteTo(outData0, outData1),
		// flow the jsons through the pipes and keep metadata
		state.Consume(),
	)

	// run the pipeline
	_, _ = pipeline.Run(context.Background(), nil)

	// print the output buffers data
	fmt.Printf("Result data 0:\n\t%s", outData0.String())
	fmt.Printf("Redult data 1:\n\t%s", outData1.String())
}
