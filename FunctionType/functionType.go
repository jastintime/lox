package functionType 

type FunctionType int


const (
	None FunctionType = iota
	Function
	Initializer
	Method
)
