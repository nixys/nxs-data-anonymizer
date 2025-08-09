// Exemple de modification pour misc/values.go
package misc

type ValueType string

const (
	ValueTypeUnknown  ValueType = "unknown"  // Corrigé la typo "unkonwn"
	ValueTypeTemplate ValueType = "template"
	ValueTypeCommand  ValueType = "command"
	ValueTypeFaker    ValueType = "faker"    // NOUVEAU TYPE AJOUTÉ
)

func ValueTypeFromString(v string) ValueType {
	switch v {
	case string(ValueTypeTemplate):
		return ValueTypeTemplate
	case string(ValueTypeCommand):
		return ValueTypeCommand
	case string(ValueTypeFaker):     // NOUVEAU CAS AJOUTÉ
		return ValueTypeFaker
	default:
		return ValueTypeUnknown
	}
}

func (v ValueType) String() string {
	return string(v)
}