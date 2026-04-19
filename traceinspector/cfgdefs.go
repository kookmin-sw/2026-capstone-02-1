package traceinspector

import (
	"encoding/json"
	"fmt"
	"strings"
	"traceinspector/imp"
)

type node_types string

type NodeID int
type EdgeID int

type CFGNodeClass interface {
	is_CFGNodeClass()
	To_mermaid() string
}

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGNodeLocation struct {
	Function_name imp.ImpFunctionName
	Id            NodeID
}

func create_empty_node_location() CFGNodeLocation {
	return CFGNodeLocation{}
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

type CFGEdgeLocation struct {
	Function_name imp.ImpFunctionName
	Id            EdgeID
}

type CFGCondEdge struct {
	Id               CFGEdgeLocation
	From_node_id     NodeID
	To_true_node_id  NodeID
	To_false_node_id NodeID
}

type CFGEdge struct {
	Id           CFGEdgeLocation
	From_node_id NodeID
	To_node_id   NodeID
	Label        string
}

type CFGEdgeClass interface {
	is_CFGEdgeClass()
}

func (node *CFGEdge) is_CFGEdgeClass() {}

func (node *CFGEdge) To_mermaid() string {
	if node.Label == "" {
		return fmt.Sprintf("%d --> %d", node.From_node_id, node.To_node_id)
	} else {
		return fmt.Sprintf("%d -- %s --> %d", node.From_node_id, node.Label, node.To_node_id)
	}
}

func (node *CFGCondEdge) is_CFGEdgeClass() {}

type CFGGraph struct {
	Entry_node    NodeID
	Node_map      map[NodeID]CFGNodeClass   // Map from node ID to node obj
	Edge_map_from map[NodeID]CFGEdgeClass   // map from node ID to outgoing edge objs
	Edge_map_to   map[NodeID][]CFGEdgeClass // map from node ID to incoming edge objs
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
	for _, edge_type := range m.Edge_map_from {
		switch edge := edge_type.(type) {
		case *CFGEdge:
			repr.Edges = append(repr.Edges, *edge)
		case *CFGCondEdge:
			if edge.To_true_node_id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_id: edge.From_node_id, To_node_id: edge.To_true_node_id, Label: "True"})
			}
			if edge.To_false_node_id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_id: edge.From_node_id, To_node_id: edge.To_false_node_id, Label: "False"})
			}
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
	for _, edge_type := range m.Edge_map_from {
		switch edge := edge_type.(type) {
		case *CFGEdge:
			repr.Edges = append(repr.Edges, *edge)
		case *CFGCondEdge:
			if edge.To_true_node_id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_id: edge.From_node_id, To_node_id: edge.To_true_node_id, Label: "True"})
			}
			if edge.To_false_node_id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_id: edge.From_node_id, To_node_id: edge.To_false_node_id, Label: "False"})
			}
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
