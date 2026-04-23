package traceinspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"traceinspector/imp"
)

type AnalyzerOutputType string

const (
	AnalyzerOutput_error       AnalyzerOutputType = "error"
	AnalyzerOutput_update_node AnalyzerOutputType = "update_node"
	AnalyzerOutput_information AnalyzerOutputType = "info"
)

type AnalyzerOutput struct {
	Type          AnalyzerOutputType
	Function_name imp.ImpFunctionName
	Node_id       NodeID
	Line_num      int
	Msg           string
}

func write_error(node_location CFGNodeLocation, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: AnalyzerOutput_error, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}

func write_update_node(node_location CFGNodeLocation, label string) {
	res := AnalyzerOutput{Type: AnalyzerOutput_update_node, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: label}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(res)
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}

func write_error_linenum(line_num int, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: AnalyzerOutput_error, Line_num: line_num, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}
