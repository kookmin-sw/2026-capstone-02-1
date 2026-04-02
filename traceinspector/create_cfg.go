package traceinspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

// if node id is leq 0, then the node doesn't exist
type CFGGraphCreator struct {
	fset                    *token.FileSet
	cfg_graph               *CFGGraph
	next_node_index         int // the next available node id
	next_edge_index         int // the next available edge id
	prev_node_index         int // the id of the last created node that should be the origin of an edge, if > 0
	prev_node_index_if_else int // additionally used to denote else statement block node id
}

func (graphcreator *CFGGraphCreator) create_cfg_node(node ast.Node, node_type node_types) int {
	current_node_index := graphcreator.next_node_index
	graphcreator.cfg_graph.Nodes = append(graphcreator.cfg_graph.Nodes, CFGNode{Id: current_node_index, Code: nodeString(node), Node_type: node_type, Line_num: graphcreator.fset.Position(node.Pos()).Line})
	graphcreator.next_node_index++
	return current_node_index
}

func (graphchreator *CFGGraphCreator) create_cfg_edge(to_id int, label string) {
	if graphchreator.prev_node_index > 0 {
		graphchreator.cfg_graph.Edges = append(graphchreator.cfg_graph.Edges, CFGEdge{Id: graphchreator.next_edge_index, From_node_id: graphchreator.prev_node_index, To_node_id: to_id, Label: label})
		graphchreator.next_edge_index++
	}
	if graphchreator.prev_node_index_if_else > 0 {
		graphchreator.cfg_graph.Edges = append(graphchreator.cfg_graph.Edges, CFGEdge{Id: graphchreator.next_edge_index, From_node_id: graphchreator.prev_node_index_if_else, To_node_id: to_id, Label: label})
		graphchreator.next_edge_index++
	}

}

// Return an AST node's code string
func nodeString(n ast.Node) string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	format.Node(&buf, fset, n)
	return buf.String()
}

func (graphcreator *CFGGraphCreator) create_cfg_method(n ast.Node) int {
	if n == nil {
		return 0
	}
	switch node := n.(type) {
	case *ast.BlockStmt:
		if len(node.List) > 0 {
			for _, subnode := range node.List {
				graphcreator.create_cfg_method(subnode)
			}
		}
		return graphcreator.prev_node_index
	case *ast.IfStmt:
		if_node_id := graphcreator.create_cfg_node(node.Cond, node_cond)
		graphcreator.create_cfg_edge(if_node_id, "")
		graphcreator.prev_node_index = if_node_id
		if node.Body != nil {
			// true body
			for index, body_node := range node.Body.List {
				if index == 0 {
					body_node_id := graphcreator.create_cfg_node(body_node, node_basic)
					graphcreator.create_cfg_edge(body_node_id, "true")
					graphcreator.prev_node_index = body_node_id
				} else {
					graphcreator.prev_node_index = graphcreator.create_cfg_method(body_node)
				}
				graphcreator.prev_node_index_if_else = graphcreator.prev_node_index
			}
		}
		if node.Else != nil {
			switch else_node := node.Else.(type) {
			case *ast.BlockStmt:
				// else body
				for index, body_node := range else_node.List {
					if index == 0 {
						body_node_id := graphcreator.create_cfg_node(body_node, node_basic)
						graphcreator.create_cfg_edge(body_node_id, "true")
						graphcreator.prev_node_index = body_node_id
					} else {
						graphcreator.prev_node_index = graphcreator.create_cfg_method(body_node)
					}
				}
			default:
				graphcreator.prev_node_index = if_node_id
				else_node_id := graphcreator.create_cfg_node(node.Else, node_basic)
				graphcreator.create_cfg_edge(else_node_id, "false")
			}
		}
		return graphcreator.prev_node_index
	default:
		node_id := graphcreator.create_cfg_node(n, node_basic)
		graphcreator.create_cfg_edge(node_id, "")
		graphcreator.prev_node_index = node_id
		return node_id
	}
}

// create and print the cfg into json
func Print_cfg(file *ast.File, fset *token.FileSet) {
	var func_cfg_map map[string]*CFGGraph = make(map[string]*CFGGraph)
	for _, decls := range file.Decls {
		switch decl_node := decls.(type) {
		case *ast.FuncDecl:
			{
				func_cfg_map[decl_node.Name.Name] = &CFGGraph{}
				cfg_creator := CFGGraphCreator{fset: fset, cfg_graph: func_cfg_map[decl_node.Name.Name], next_node_index: 1}
				cfg_creator.create_cfg_method(decl_node.Body)
			}
		}
	}
	// result, _ := json.Marshal(func_cfg_map)
	result, _ := json.MarshalIndent(func_cfg_map, "", "    ") // return formatted
	fmt.Println(string(result))
}
