package misc

type SecurityPolicyTablesType string

const (
	SecurityPolicyTablesUnknown SecurityPolicyTablesType = "unknown"
	SecurityPolicyTablesPass    SecurityPolicyTablesType = "pass"
	SecurityPolicyTablesSkip    SecurityPolicyTablesType = "skip"
)

func (v SecurityPolicyTablesType) String() string {
	return string(v)
}

func SecurityPolicyTablesTypeFromString(v string) SecurityPolicyTablesType {
	switch v {
	case string(SecurityPolicyTablesPass):
		return SecurityPolicyTablesPass
	case string(SecurityPolicyTablesSkip):
		return SecurityPolicyTablesSkip
	default:
		return SecurityPolicyTablesUnknown
	}
}

type SecurityPolicyColumnsType string

const (
	SecurityPolicyColumnsUnknown   SecurityPolicyColumnsType = "unknown"
	SecurityPolicyColumnsPass      SecurityPolicyColumnsType = "pass"
	SecurityPolicyColumnsRandomize SecurityPolicyColumnsType = "randomize"
)

func (v SecurityPolicyColumnsType) String() string {
	return string(v)
}

func SecurityPolicyColumnsTypeFromString(v string) SecurityPolicyColumnsType {
	switch v {
	case string(SecurityPolicyColumnsPass):
		return SecurityPolicyColumnsPass
	case string(SecurityPolicyColumnsRandomize):
		return SecurityPolicyColumnsRandomize
	default:
		return SecurityPolicyColumnsUnknown
	}
}
