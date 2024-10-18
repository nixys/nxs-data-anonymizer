package pgsql_anonymize

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhSecurityCopy(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = append(uctx.security.tmpBuf, token...)
	uctx.insertIntoBuf = nil

	return deferred, nil
}

func dhCopyValuesEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if uctx.insertIntoBuf != nil {
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

	clmns := make(map[string]string)

	ss := bytes.Split(deferred, []byte{'\n'})

	for _, v := range ss {

		s := strings.TrimSuffix(strings.TrimSpace(string(v)), ",")

		if len(s) > 0 {

			u := strings.SplitN(s, " ", 2)

			// If column type does not specified within the dump
			if len(u) < 2 {
				clmns[u[0]] = ""
			} else {
				clmns[u[0]] = u[1]
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

	uctx.insertIntoBuf = append(uctx.security.tmpBuf, append(deferred, token...)...)

	uctx.security.isSkip = false
	uctx.security.tmpBuf = []byte{}

	return []byte{}, nil
}

func dhFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	fname := bytes.Trim(bytes.TrimSpace(deferred), "\"")

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.filter.ColumnAdd(
		string(fname),
		uctx.tables[uctx.filter.TableNameGet()][string(fname)],
	)

	uctx.insertIntoBuf = append(uctx.insertIntoBuf, append(deferred, token...)...)

	return []byte{}, nil
}

func dhTableCopyTail(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.insertIntoBuf = append(uctx.insertIntoBuf, append(deferred, token...)...)

	return []byte{}, nil
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

	b := rowDataGen(uctx.filter)
	if b == nil {
		return []byte{}, nil
	} else {
		if uctx.insertIntoBuf != nil {
			b = append(uctx.insertIntoBuf, b...)
			uctx.insertIntoBuf = nil
		}
	}

	return b, nil
}

func rowDataGen(filter *relfilter.Filter) []byte {

	var out string

	row := filter.ValuePop()
	if row.Values == nil {
		return nil
	}

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
	if uctx.security.tablesPolicy != misc.SecurityPolicyTablesSkip {
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
