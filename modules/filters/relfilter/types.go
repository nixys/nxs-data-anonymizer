package relfilter

type ValueType int

const (
	ValueTypeByte ValueType = iota
	ValueTypeString
	ValueTypeBinary
	ValueTypeInt
	ValueTypeFloat
	ValueTypeNULL
)

func (vt ValueType) String() string {
	return [...]string{"byte", "string", "binary", "int", "float", "null"}[vt]
}
