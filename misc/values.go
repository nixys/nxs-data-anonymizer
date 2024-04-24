package misc

type ValueType string

const (
	ValueTypeUnknown  ValueType = "unkonwn"
	ValueTypeTemplate ValueType = "template"
	ValueTypeCommand  ValueType = "command"
)

func ValueTypeFromString(v string) ValueType {
	switch v {
	case string(ValueTypeTemplate):
		return ValueTypeTemplate
	case string(ValueTypeCommand):
		return ValueTypeCommand
	default:
		return ValueTypeUnknown
	}
}

func (v ValueType) String() string {
	return string(v)
}
