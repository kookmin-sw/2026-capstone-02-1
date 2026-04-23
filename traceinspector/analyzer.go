package traceinspector

import (
	"traceinspector/domain"
	"traceinspector/imp"
)

// This is the main struct for driving abstract interpretation.
type AbstractAnalyzer[IntDom domain.IntegerDomain[IntDom], ArrDom ArrayDomain[IntDom, ArrDom]] struct {
	Function_cfgs         FunctionCFGMap
	Function_defs         imp.ImpFunctionMap
	Create_semantics_func func(*FunctionAbstractMem[IntDom, ArrDom], imp.ImpFunctionName) AbstractSemantics[IntDom, ArrDom]
	function_pre_mem_map  map[imp.ImpFunctionName]*FunctionAbstractMem[IntDom, ArrDom] // map from function name to pre-states
}

func (analyzer *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Start_analysis(function_name imp.ImpFunctionName) {
	analyzer.function_pre_mem_map = make(map[imp.ImpFunctionName]*FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl])
	analyzer.function_pre_mem_map[function_name] = &FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl]{}
	analyzer.function_pre_mem_map[function_name].Initialize(function_name, analyzer.Function_cfgs[function_name])

	semantics := analyzer.Create_semantics_func(analyzer.function_pre_mem_map[function_name], function_name)

	initial_state := AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: analyzer.Function_cfgs[function_name].Entry_node, abstract_mem: make(AbstractNodeMem[IntDomainImpl, ArrayDomainImpl])}
	worklist := []AbstractState[IntDomainImpl, ArrayDomainImpl]{initial_state}
	for len(worklist) > 0 {
		front_val := worklist[0]
		worklist = worklist[1:]
		// fmt.Println("Process state", front_val)
		for _, val := range semantics.Step(front_val) {
			worklist = append(worklist, val)
		}
	}
	// fmt.Println("Final mem", analyzer.function_pre_mem_map[function_name])
}

func Test(func_cfg_map FunctionCFGMap, func_name imp.ImpFunctionName, func_info_map imp.ImpFunctionMap) {
	create_sem := func(func_mem *FunctionAbstractMem[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]], func_name imp.ImpFunctionName) AbstractSemantics[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]] {
		return &ImpFunctionInterpreter[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{
			func_cfg_map:        func_cfg_map,
			func_name:           func_name,
			func_info_map:       func_info_map,
			abstract_mem:        func_mem,
			intdomain_default:   domain.IntervalDomain{},
			booldomain_default:  domain.BoolDomain{},
			arraydomain_default: ArraySummaryDomain[domain.IntervalDomain]{},
		}
	}
	g := AbstractAnalyzer[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{Function_cfgs: func_cfg_map, Function_defs: func_info_map, Create_semantics_func: create_sem}
	g.Start_analysis("main")
}
