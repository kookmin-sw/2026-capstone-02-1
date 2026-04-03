package traceinspector

import (
	"encoding/json"
	"fmt"
)

type node_types string

type CFGComponent interface {
	to_json() []byte
}

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGNode struct {
	Id        int
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGEdge struct {
	Id           int
	From_node_id int
	To_node_id   int
	Label        string
}

func (node *CFGNode) to_json() []byte {
	out, err := json.Marshal(node)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal CFG Node %s\n", node.Code))
	}
	return out
}

func (node *CFGEdge) to_json() []byte {
	out, err := json.Marshal(node)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal CFG Edge %s\n", node.Label))
	}
	return out
}

type CFGGraph struct {
	Nodes []CFGNode
	Edges []CFGEdge
}

func (cfg_graph *CFGGraph) to_json() []byte {
	out, err := json.Marshal(cfg_graph)
	if err != nil {
		panic("Failed to marshal CFG graph")
	}
	return out
}
