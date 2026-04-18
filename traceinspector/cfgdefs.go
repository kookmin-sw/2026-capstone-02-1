package traceinspector

import (
	"encoding/json"
	"fmt"
	"strings"
	"traceinspector/imp"
)

type node_types string

type CFGNodeClass interface {
	is_CFGNodeClass()
	To_mermaid() string
}

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGNodeLocation struct {
	Function_name string
	Id            int
}

type CFGNode struct {
	Ast       imp.Stmt `json:"-"`
	Id        CFGNodeLocation
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGCondNode struct {
	Ast       imp.Expr `json:"-"`
	Id        CFGNodeLocation
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGEdge struct {
	Id           CFGNodeLocation
	From_node_id int
	To_node_id   int
	Label        string
}

func (node *CFGNode) is_CFGNodeClass() {}

func (node *CFGNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, node.Code)
	}
	return ""
}

func (node *CFGCondNode) is_CFGNodeClass() {}

func (node *CFGCondNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, node.Code)
	}
	return ""
}

func (node *CFGEdge) To_mermaid() string {
	if node.Label == "" {
		return fmt.Sprintf("%d --> %d", node.From_node_id, node.To_node_id)
	} else {
		return fmt.Sprintf("%d -- %s --> %d", node.From_node_id, node.Label, node.To_node_id)
	}
}

type CFGGraph struct {
	Node_map      map[int]CFGNodeClass // Map from node ID to node obj
	Edge_map_from map[int][]*CFGEdge   // map from node ID to outgoing edge objs
	Edge_map_to   map[int][]*CFGEdge   // map from node ID to incoming edge objs
}

func (m CFGGraph) MarshalJSON() ([]byte, error) {
	type CFGGraphRepr struct {
		Nodes []CFGNodeClass
		Edges []CFGEdge
	}
	repr := CFGGraphRepr{}
	for _, v := range m.Node_map {
		repr.Nodes = append(repr.Nodes, v)
	}
	for _, edges := range m.Edge_map_from {
		for _, edge := range edges {
			repr.Edges = append(repr.Edges, *edge)
		}
	}
	return json.Marshal(repr)
}

func (m CFGGraph) To_mermaid() string {
	type CFGGraphRepr struct {
		Nodes []CFGNodeClass
		Edges []CFGEdge
	}
	repr := CFGGraphRepr{}
	out := strings.Builder{}
	for _, v := range m.Node_map {
		repr.Nodes = append(repr.Nodes, v)
	}
	for _, edges := range m.Edge_map_from {
		for _, edge := range edges {
			repr.Edges = append(repr.Edges, *edge)
		}
	}
	for _, node := range repr.Nodes {
		out.WriteString(node.To_mermaid() + "\n")
	}

	for _, edge := range repr.Edges {
		out.WriteString(edge.To_mermaid() + "\n")
	}
	return out.String()
}
