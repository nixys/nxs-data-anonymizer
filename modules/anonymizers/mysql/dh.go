package mysql_anonymize

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhSecurityCreateTable(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = append(uctx.security.tmpBuf, token...)

	return deferred, nil
}

func dhSecurityCreateTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = append(uctx.security.tmpBuf, deferred...)
	uctx.security.tmpBuf = append(uctx.security.tmpBuf, token...)

	return []byte{}, nil
}

func dhSecurityNil(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	return append(deferred, token...), nil
}

func dhCreateTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	tn := string(deferred)

	// Check table pass through security rules
	if !securityPolicyCheck(uctx, tn) {

		// If not: table will be skipped from result dump

		uctx.security.isSkip = true
		uctx.security.tmpBuf = []byte{}
		return []byte{}, nil
	}

	uctx.filter.TableCreate(tn)

	d := append(uctx.security.tmpBuf, append(deferred, token...)...)

	uctx.security.isSkip = false
	uctx.security.tmpBuf = []byte{}

	return d, nil
}

func dhCreateTableFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.column.name = string(deferred)

	return append(deferred, token...), nil
}

func dhCreateTableColumnTypeAdd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	for k, v := range typeKeys {
		if k == "generated" {
			if k == string(token) || strings.ToUpper(k) == string(token) {
				uctx.column.isSkip = true
				break
			}
		} else {
			if k == string(token) || strings.ToUpper(k) == string(token) {
				uctx.column.columnType = v
				break
			}
		}
	}

	if uctx.column.columnType == "" {
		fmt.Println("token:", token)
	}

	return append(deferred, token...), nil
}

func dhCreateTableColumnAdd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if uctx.column.isSkip == false {
		uctx.filter.ColumnAdd(uctx.column.name, uctx.column.columnType)
	}

	uctx.column = userColumnCtx{}

	return append(deferred, token...), nil
}

func dhInsertIntoTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	// Check insert into table name
	if bytes.Compare([]byte(uctx.filter.TableNameGet()), deferred) != 0 {
		return append(deferred, token...), fmt.Errorf("`create` and `insert into` table names are mismatch (create table: '%s', insert into table: '%s')", uctx.filter.TableNameGet(), string(deferred))
	}

	return append(deferred, token...), nil
}

func dhCreateTableValues(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if bytes.Compare(deferred, []byte("NULL")) == 0 {
		uctx.filter.ValueAdd(nil)
	} else {
		uctx.filter.ValueAdd(deferred)
	}

	return []byte{}, nil
}

func dhCreateTableValuesString(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.filter.ValueAdd(deferred)

	return []byte{}, nil
}

func dhCreateTableValuesEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	if bytes.Compare(deferred, []byte("NULL")) == 0 {
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

func dhCreateTableValuesStringEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
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
			out += ","
		}

		if v.V == nil {
			out += "NULL"
		} else {
			switch filter.ColumnTypeGet(i) {
			case relfilter.ColumnTypeString:
				out += fmt.Sprintf("'%s'", v.V)
			case relfilter.ColumnTypeBinary:
				out += fmt.Sprintf("_binary '%s'", v.V)
			default:
				out += fmt.Sprintf("%s", v.V)
			}
		}
	}

	return []byte(fmt.Sprintf("(%s)", out))
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
