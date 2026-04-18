package traceinspector

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"strings"
	"traceinspector/imp"
)

// if node id is leq 0, then the node doesn't exist
type CFGGraphCreator struct {
	func_name         string // name of the function
	fset              *token.FileSet
	Cfg_graph         *CFGGraph
	next_node_index   int          // the next available node id
	next_edge_index   int          // the next available edge id
	cfg_context_stack []CFGContext // stack holding the graph context

}

type CFGContext interface {
	isCFGContext()
}

type CFGLoopContext struct {
	head_node_id int // node ID of the loop head(condition node)
	exit_node_id int // node ID of the node after the loop
}

func (CFGLoopContext) isCFGContext() {}

type CFGBranchContext struct {
	exit_node_id int // node ID of the node after the branch(join node)
}

func (CFGBranchContext) isCFGContext() {}

// Return the topmost loop context
func (creator *CFGGraphCreator) get_top_loop_context() *CFGLoopContext {
	for stack_index := len(creator.cfg_context_stack) - 1; stack_index >= 0; stack_index-- {
		loop_context, is_loop_context := creator.cfg_context_stack[stack_index].(CFGLoopContext)
		if is_loop_context {
			return &loop_context
		}
	}
	return nil
}

// Return the next stmt node ID to evaluate(link to), given the current get_top_context_destination
// If in a branch, it's the stmt after the branch
// If in a loop, it's the loop head(condition node)
func (creator *CFGGraphCreator) get_top_context_destination() int {
	for stack_index := len(creator.cfg_context_stack) - 1; stack_index >= 0; stack_index-- {
		switch ctx := creator.cfg_context_stack[stack_index].(type) {
		case CFGLoopContext:
			if ctx.exit_node_id > 0 {
				return ctx.head_node_id
			}
		case CFGBranchContext:
			if ctx.exit_node_id > 0 {
				return ctx.exit_node_id
			}
		}
	}
	return 0
}

func (creator *CFGGraphCreator) push_branch_context(cond_node_id int, exit_node_id int) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, CFGBranchContext{exit_node_id})
}

func (creator *CFGGraphCreator) push_loop_context(cond_node_id int, exit_node_id int) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, CFGLoopContext{cond_node_id, exit_node_id})
}

func (creator *CFGGraphCreator) pop_context() {
	creator.cfg_context_stack = creator.cfg_context_stack[:len(creator.cfg_context_stack)-1]
}

func (graphcreator *CFGGraphCreator) create_cfg_node(imp_ast imp.Stmt, line_num int) int {
	current_node_index := graphcreator.next_node_index
	escaped_code := strings.ReplaceAll(fmt.Sprintf("%s", imp_ast), "\"", "#34;")
	graphcreator.Cfg_graph.Node_map[current_node_index] = &CFGNode{Ast: imp_ast, Id: CFGNodeLocation{graphcreator.func_name, current_node_index}, Code: escaped_code, Node_type: node_basic, Line_num: line_num}
	graphcreator.next_node_index++
	return current_node_index
}

func (graphcreator *CFGGraphCreator) create_cfg_cond_node(imp_ast imp.Expr, line_num int) int {
	current_node_index := graphcreator.next_node_index
	escaped_code := strings.ReplaceAll(fmt.Sprintf("%s", imp_ast), "\"", "#34;")
	graphcreator.Cfg_graph.Node_map[current_node_index] = &CFGCondNode{Ast: imp_ast, Id: CFGNodeLocation{graphcreator.func_name, current_node_index}, Code: escaped_code, Node_type: node_cond, Line_num: line_num}
	graphcreator.next_node_index++
	return current_node_index
}

func (graphcreator *CFGGraphCreator) create_cfg_edge(from_id int, to_id int, label string) {
	if from_id > 0 && to_id > 0 {
		escaped_label := strings.ReplaceAll(label, "\"", "#34;")
		edge := CFGEdge{Id: CFGNodeLocation{graphcreator.func_name, graphcreator.next_edge_index}, From_node_id: from_id, To_node_id: to_id, Label: escaped_label}
		graphcreator.Cfg_graph.Edge_map_from[from_id] = append(graphcreator.Cfg_graph.Edge_map_from[from_id], &edge)
		graphcreator.Cfg_graph.Edge_map_to[to_id] = append(graphcreator.Cfg_graph.Edge_map_to[to_id], &edge)
		graphcreator.next_edge_index++
	}
}

// The driver function for creating the CFG graph. stmt is the current statement node.
// linkback, if not 0, equals the node id that an edge should be created from the current node to the linkback ID
func (graphcreator *CFGGraphCreator) create_cfg_method(stmts []imp.Stmt) int {
	if len(stmts) == 0 {
		return 0
	}
	next_node_id := graphcreator.create_cfg_method(stmts[1:]) // slice[1:] returns empty slice for len 1 slice
	if next_node_id == 0 {
		// If there's no remaining statement, the next destination depends on context
		next_node_id = graphcreator.get_top_context_destination()
	}
	var created_node_id int = 0
	switch stmt_ty := stmts[0].(type) {
	case *imp.IfElseStmt:
		cond_node_id := graphcreator.create_cfg_cond_node(stmt_ty.Cond, stmt_ty.GetLineNum())

		graphcreator.push_branch_context(cond_node_id, next_node_id)

		// the node ID of the starting node in true stmt flow
		true_node_id := graphcreator.create_cfg_method(stmt_ty.True_stmt)
		if true_node_id == 0 {
			// true stmt empty, next destination is context-dependent
			true_node_id = next_node_id
		}

		// the node ID of the starting node in the false stmt flow
		false_node_id := graphcreator.create_cfg_method(stmt_ty.False_stmt)
		if false_node_id == 0 {
			// false stmt empty, next destination is context dependent
			false_node_id = next_node_id
		}

		// create edges from cond to true/false start node
		graphcreator.create_cfg_edge(cond_node_id, true_node_id, "True")
		graphcreator.create_cfg_edge(cond_node_id, false_node_id, "False")

		graphcreator.pop_context()

		created_node_id = cond_node_id

	case *imp.WhileStmt:
		cond_node_id := graphcreator.create_cfg_cond_node(stmt_ty.Cond, stmt_ty.GetLineNum())

		graphcreator.push_loop_context(cond_node_id, next_node_id)
		body_node_id := graphcreator.create_cfg_method(stmt_ty.Body_stmt)
		graphcreator.create_cfg_edge(cond_node_id, body_node_id, "True")
		graphcreator.create_cfg_edge(cond_node_id, next_node_id, "False")
		graphcreator.pop_context()

		created_node_id = cond_node_id

	case *imp.BreakStmt:
		created_node_id = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		ctx := graphcreator.get_top_loop_context()
		// link to loop exit
		graphcreator.create_cfg_edge(created_node_id, ctx.exit_node_id, "")

	case *imp.ContinueStmt:
		created_node_id = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		ctx := graphcreator.get_top_loop_context()
		// link to loop head
		graphcreator.create_cfg_edge(created_node_id, ctx.head_node_id, "")

	case *imp.ReturnStmt:
		created_node_id = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		// finish generation
	default:
		created_node_id = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		graphcreator.create_cfg_edge(created_node_id, next_node_id, "")

	}
	return created_node_id
}

// create and print the cfg into json
func Create_cfg(functions map[string]imp.ImpFunction) map[string]*CFGGraph {
	var func_cfg_map map[string]*CFGGraph = make(map[string]*CFGGraph)
	for fun_name, fun := range functions {
		func_cfg_map[fun_name] = &CFGGraph{Node_map: make(map[int]CFGNodeClass), Edge_map_from: map[int][]*CFGEdge{}, Edge_map_to: map[int][]*CFGEdge{}}
		cfg_creator := CFGGraphCreator{func_name: fun_name, Cfg_graph: func_cfg_map[fun_name], next_node_index: 1}
		cfg_creator.create_cfg_method(fun.Body)
	}
	return func_cfg_map
}

func Print_cfg_map_json(cfgs map[string]*CFGGraph) {
	// result, _ := json.Marshal(func_cfg_map)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	enc.Encode(cfgs)
}

func Print_mermaid(cfg *CFGGraph) {
	// fmt.Println("```")
	fmt.Println("flowchart TD")
	fmt.Println(cfg.To_mermaid())
	// fmt.Println("```")
}
