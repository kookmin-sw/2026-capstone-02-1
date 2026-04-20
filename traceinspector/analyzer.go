package traceinspector

import (
	"traceinspector/domain"
	"traceinspector/imp"
)

// This is the main struct for initializing abstract interpretation.
//
// abstract_semantics: The interface implementing the abstract step relation
// function_mem_map:
type AbstractAnalyzer[IntDom domain.AbstractDomain[IntDom], BoolDom domain.AbstractDomain[BoolDom], ArrDom domain.AbstractDomain[ArrDom]] struct {
	abstract_semantics AbstractSemantics[IntDom, BoolDom, ArrDom]
	function_mem_map   map[imp.ImpFunctionName]*FunctionAbstractMem[IntDom, BoolDom, ArrDom]
	function_cfgs      FunctionCFGMap
	function_defs      imp.ImpFunctionMap
}

func (analyzer *AbstractAnalyzer[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) Start_analysis(function_name imp.ImpFunctionName) {
	analyzer.function_mem_map = make(map[imp.ImpFunctionName]*FunctionAbstractMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl])
	analyzer.function_mem_map[function_name] = &FunctionAbstractMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]{}
	analyzer.function_mem_map[function_name].Initialize(function_name)

	// initial_state := AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]{}
}

// func test() {
// 	semantics := ImpAbstractSemantics[IntervalDomain, BoolDomain, ArraySummaryDomain[IntervalDomain]]{}
// 	g := AbstractAnalyzer[IntervalDomain, BoolDomain, ArraySummaryDomain[IntervalDomain]]{abstract_semantics: &semantics}

// }
