package traceinspector

import (
	"encoding/json"
	"os"
	"traceinspector/imp"
)

type AnalyzerOutputType string

const (
	AnalyzerOutput_error       AnalyzerOutputType = "error"
	AnalyzerOutput_update_node AnalyzerOutputType = "update_node"
	AnalyzerOutput_info        AnalyzerOutputType = "info"
	AnalyzerOutput_warning     AnalyzerOutputType = "warning"
)

type AnalyzerOutputHandler struct {
	Debugs []AnalyzerOutput
}

type AnalyzerOutput struct {
	Type          AnalyzerOutputType
	Function_name imp.ImpFunctionName
	Node_id       NodeID
	Node_state    string
	Msg           string
}

func (ao *AnalyzerOutputHandler) Print() {
	// buf := &bytes.Buffer{}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	enc.Encode(ao)
	// out := &bytes.Buffer{}
	// json.Compact(out, buf.Bytes())
	// fmt.Println(out.String())
}

func (ao *AnalyzerOutputHandler) write_info(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_info, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
}

func (ao *AnalyzerOutputHandler) write_warning(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_warning, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
}

func (ao *AnalyzerOutputHandler) write_error(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_error, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	ao.Print()
	os.Exit(1)
}

func (ao *AnalyzerOutputHandler) write_update_node_state(node_location CFGNodeLocation, state_str string, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_update_node, Function_name: node_location.Function_name, Node_id: node_location.Id, Node_state: state_str, Msg: msg})
}
