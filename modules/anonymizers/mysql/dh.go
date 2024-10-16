package mysql_anonymize

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/nixys/nxs-data-anonymizer/misc"
)

func dhSecurityInsertInto(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = token

	return deferred, nil
}

func dhSecurityInsertIntoTableNameSearch(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.security.tmpBuf = append(uctx.security.tmpBuf, deferred...)
	uctx.security.tmpBuf = append(uctx.security.tmpBuf, token...)

	return []byte{}, nil
}

func dhSecurityInsertIntoValues(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.insertIntoBuf = append(uctx.insertIntoBuf, deferred...)
	uctx.insertIntoBuf = append(uctx.insertIntoBuf, token...)

	return []byte{}, nil
}

func dhSecurityInsertIntoValueSearch(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
		return []byte{}, nil
	}

	uctx.insertIntoBuf = append(uctx.insertIntoBuf, deferred...)

	return []byte{}, nil
}

func dhSecurityValuesEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

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

	uctx := usrCtx.(*userCtx)
	uctx.filter.TableCreate(string(deferred))

	return append(deferred, token...), nil
}

func dhCreateTableFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)
	uctx.columnName = string(deferred)

	return append(deferred, token...), nil
}

// checkGenerated returns true when specified type is `generated`
// See: https://dev.mysql.com/blog-archive/generated-columns-in-mysql-5-7-5 for details
func checkGenerated(t []byte) bool {
	if bytes.Contains(t, []byte{'A', 'S'}) == true {
		b, _ := regexp.Match("^([A-Z]+)((\\([0-9]+\\) )| )(GENERATED ALWAYS AS|AS)", t)
		return b
	}
	return false
}

func dhCreateTableColumnAdd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	traw := bytes.TrimSpace(deferred)
	trawUpper := bytes.ToUpper(traw)

	if checkGenerated(trawUpper) == false {

		t, b := uctx.tables[uctx.filter.TableNameGet()]
		if !b {
			t = make(map[string]columnType)
		}
		t[uctx.columnName] = columnTypeNone

		for _, ot := range uctx.optKinds {
			if ot.r.Match(trawUpper) == true {
				t[uctx.columnName] = ot.t
				break
			}
		}

		uctx.tables[uctx.filter.TableNameGet()] = t
		uctx.filter.ColumnAdd(uctx.columnName, string(traw))
	}

	uctx.columnName = ""

	return append(deferred, token...), nil
}

func dhInsertIntoTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	tn := string(deferred)

	// Check table pass through security rules
	if !securityPolicyCheck(uctx, tn) {

		// If not: table will be skipped from result dump

		uctx.security.isSkip = true
		uctx.security.tmpBuf = []byte{}
		return []byte{}, nil
	}

	uctx.insertIntoBuf = append(uctx.security.tmpBuf, append(deferred, token...)...)

	uctx.security.isSkip = false
	uctx.security.tmpBuf = []byte{}

	// Check insert into table name
	if tn != uctx.filter.TableNameGet() {
		return []byte{}, fmt.Errorf("`create` and `insert into` table names are mismatch (create table: '%s', insert into table: '%s')", uctx.filter.TableNameGet(), tn)
	}

	return []byte{}, nil
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

	b := rowDataGen(uctx)
	if b == nil {
		return []byte{}, nil
	} else {
		if uctx.insertIntoBuf != nil {
			b = append(uctx.insertIntoBuf, b...)
			uctx.insertIntoBuf = nil
		} else {
			b = append([]byte{','}, b...)
		}
	}

	return b, nil
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

	b := rowDataGen(uctx)
	if b == nil {
		return []byte{}, nil
	} else {
		if uctx.insertIntoBuf != nil {
			b = append(uctx.insertIntoBuf, b...)
			uctx.insertIntoBuf = nil
		} else {
			b = append([]byte{','}, b...)
		}
	}

	return b, nil
}

func rowDataGen(uctx *userCtx) []byte {

	var out string

	row := uctx.filter.ValuePop()
	if row.Values == nil {
		return nil
	}

	for i, v := range row.Values {

		if i > 0 {
			out += ","
		}

		if v.V == nil {
			out += "NULL"
		} else {
			switch uctx.tables[uctx.filter.TableNameGet()][uctx.filter.ColumnGetName(i)] {
			case columnTypeString:
				out += fmt.Sprintf("'%s'", v.V)
			case columnTypeBinary:
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
