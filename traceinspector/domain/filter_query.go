package domain

type FilterQuery int

const (
	FilterQuery_Eq FilterQuery = iota
	FilterQuery_Neq
	FilterQuery_Leq
	FilterQuery_Geq
)
