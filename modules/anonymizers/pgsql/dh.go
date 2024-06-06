package pgsql_anonymize

import (
	"bytes"
	"fmt"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhSecurityCopy(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = append(uctx.security.tmpBuf, token...)

	return deferred, nil
}

func dhSecurityNil(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	return append(deferred, token...), nil
}

func dhCreateTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	tname := string(bytes.TrimSpace(deferred))

	uctx := usrCtx.(*userCtx)
	uctx.tn = &tname

	return append(deferred, token...), nil
}

func dhCreateTableDesc(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	clmns := make(map[string]relfilter.ColumnType)

	ss := bytes.Split(deferred, []byte{'\n'})

	for _, s := range ss {

		s = bytes.TrimSuffix(bytes.TrimSpace(s), []byte{','})

		if len(s) > 0 {

			u := bytes.SplitN(s, []byte{' '}, 3)

			// If column type does not specified within the dump
			if len(u) < 2 {
				clmns[string(u[0])] = relfilter.ColumnTypeNone
			} else {
				clmns[string(u[0])] = columnType(string(u[1]))
			}
		}
	}

	uctx.tables[*uctx.tn] = clmns
	uctx.tn = nil

	return append(deferred, token...), nil
}

func dhTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	tname := string(bytes.TrimSpace(deferred))

	if !securityPolicyCheck(uctx, tname) {

		// If not: table will be skipped from result dump

		uctx.security.isSkip = true
		uctx.security.tmpBuf = []byte{}

		return []byte{}, nil
	}

	uctx.filter.TableCreate(tname)

	d := append(uctx.security.tmpBuf, append(deferred, token...)...)

	uctx.security.isSkip = false
	uctx.security.tmpBuf = []byte{}

	return d, nil
}

func dhFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	fname := bytes.Trim(bytes.TrimSpace(deferred), "\"")

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	t, b := uctx.tables[uctx.filter.TableNameGet()][string(fname)]
	if b == false {
		t = relfilter.ColumnTypeNone
	}

	uctx.filter.ColumnAdd(string(fname), t)

	return append(deferred, token...), nil
}

func dhValue(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if bytes.Compare(deferred, []byte("\\N")) == 0 {
		uctx.filter.ValueAdd(nil)
	} else {
		uctx.filter.ValueAdd(deferred)
	}

	return []byte{}, nil
}

func dhValueEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if bytes.Compare(deferred, []byte("\\N")) == 0 {
		uctx.filter.ValueAdd(nil)
	} else {
		uctx.filter.ValueAdd(deferred)
	}

	// Apply filter for row
	if err := uctx.filter.Apply(); err != nil {
		return []byte{}, err
	}

	return rowDataGen(uctx.filter), nil
}

func rowDataGen(filter *relfilter.Filter) []byte {

	var out string

	row := filter.ValuePop()

	for i, v := range row.Values {

		if i > 0 {
			out += "\t"
		}

		if v.V == nil {
			out += "\\N"
		} else {
			out += fmt.Sprintf("%s", v.V)
		}
	}

	return []byte(fmt.Sprintf("%s\n", out))
}

// SecurityPolicyCheck checks the table passes the security rules
// true:  pass
// false: skip
func securityPolicyCheck(uctx *userCtx, tname string) bool {

	// Continue if security policy is `skip`
	if uctx.security.tablePolicy != misc.SecurityPolicyTablesSkip {
		return true
	}

	// Check rules for specified table name
	if tr := uctx.filter.TableRulesLookup(tname); tr != nil {
		return true
	}

	// Check specified table name in exceptions
	if _, b := uctx.security.tableExceptions[tname]; b == true {
		return true
	}

	return false
}
