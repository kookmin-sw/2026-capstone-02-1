package traceinspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"traceinspector/imp"
)

type AnalyzerOutput struct {
	Type          string
	Function_name imp.ImpFunctionName
	Node_id       NodeID
	Msg           string
}

func write_error(node_location CFGNodeLocation, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: "error", Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}
