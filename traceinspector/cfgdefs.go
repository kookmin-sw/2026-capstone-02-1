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

func (loc CFGNodeLocation) String() string {
	return fmt.Sprintf("Node %d @ func %s", loc.Id, loc.Function_name)
}

func (loc CFGNodeLocation) NodeExists() bool {
	return loc.Id > 0
}

func (loc CFGNodeLocation) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(loc.Id)), nil
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
	Cond_expr    imp.Expr `json:"-"`
	Id           CFGNodeLocation
	Code         string
	Node_type    node_types
	Line_num     int
	Is_loop_head bool `json:"-"` // whether the cond is the head of the node
}

func (node *CFGNode) is_CFGNodeClass() {}

func (node *CFGNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, escape_string_mermaid(node.Code))
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, escape_string_mermaid(node.Code))
	}
	return ""
}

func (node *CFGCondNode) is_CFGNodeClass() {}

func (node *CFGCondNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, escape_string_mermaid(node.Code))
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, escape_string_mermaid(node.Code))
	}
	return ""
}

type CFGEdgeLocation struct {
	Function_name imp.ImpFunctionName
	Id            EdgeID
}

func (loc CFGEdgeLocation) String() string {
	return fmt.Sprintf("Edge %d @ func %s", loc.Id, loc.Function_name)
}

func (loc CFGEdgeLocation) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(loc.Id)), nil
}

type CFGCondEdge struct {
	Id                CFGEdgeLocation
	From_node_loc     CFGNodeLocation
	To_true_node_loc  CFGNodeLocation
	To_false_node_loc CFGNodeLocation
}

func (edge CFGCondEdge) String() string {
	return fmt.Sprintf("(%s){true : %s -> %s, false : %s -> %s}", edge.Id, edge.From_node_loc, edge.To_true_node_loc, edge.From_node_loc, edge.To_false_node_loc)
}

type CFGEdge struct {
	Id            CFGEdgeLocation
	From_node_loc CFGNodeLocation
	To_node_loc   CFGNodeLocation
	Label         string
}

func (edge CFGEdge) String() string {
	return fmt.Sprintf("(%s){%s -> %s}", edge.Id, edge.From_node_loc, edge.To_node_loc)
}

type CFGEdgeClass interface {
	is_CFGEdgeClass()
}

func (node *CFGEdge) is_CFGEdgeClass() {}

func (node *CFGEdge) To_mermaid() string {
	if node.Label == "" {
		return fmt.Sprintf("%d --> %d", node.From_node_loc.Id, node.To_node_loc.Id)
	} else {
		return fmt.Sprintf("%d -- %s --> %d", node.From_node_loc.Id, node.Label, node.To_node_loc.Id)
	}
}

func (node *CFGCondEdge) is_CFGEdgeClass() {}

type CFGGraph struct {
	Entry_node    CFGNodeLocation
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
			if edge.To_true_node_loc.Id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_loc: edge.From_node_loc, To_node_loc: edge.To_true_node_loc, Label: "True"})
			}
			if edge.To_false_node_loc.Id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_loc: edge.From_node_loc, To_node_loc: edge.To_false_node_loc, Label: "False"})
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
	out.WriteString("flowchart TD\n")
	for _, v := range m.Node_map {
		repr.Nodes = append(repr.Nodes, v)
	}
	for _, edge_type := range m.Edge_map_from {
		switch edge := edge_type.(type) {
		case *CFGEdge:
			repr.Edges = append(repr.Edges, *edge)
		case *CFGCondEdge:
			if edge.To_true_node_loc.Id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_loc: edge.From_node_loc, To_node_loc: edge.To_true_node_loc, Label: "True"})
			}
			if edge.To_false_node_loc.Id > 0 {
				repr.Edges = append(repr.Edges, CFGEdge{Id: edge.Id, From_node_loc: edge.From_node_loc, To_node_loc: edge.To_false_node_loc, Label: "False"})
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

func escape_string_mermaid(input string) string {
	input = strings.ReplaceAll(input, "`", "#96;")
	input = strings.ReplaceAll(input, "\"", "#34;")
	return input
}
