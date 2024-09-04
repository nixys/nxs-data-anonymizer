package mysql_anonymize

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

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

func dhSecurityNil(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if uctx.security.isSkip == true {
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
func checkGenerated(t string) bool {
	if strings.Contains(t, "AS") == true {
		b, _ := regexp.Match("^([A-Z]+)((\\([0-9]+\\) )| )(GENERATED ALWAYS AS|AS)", []byte(t))
		return b
	}
	return false
}

func dhCreateTableColumnAdd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	traw := strings.TrimSpace(string(deferred))
	trawUpper := strings.ToUpper(traw)

	if checkGenerated(trawUpper) == false {

		i := strings.IndexAny(strings.TrimSpace(trawUpper), " (,")
		if i != -1 {

			ct := columnTypeNone
			for k, v := range typeKeys {
				if trawUpper[0:i] == k {
					ct = v
				}
			}

			t, b := uctx.tables[uctx.filter.TableNameGet()]
			if b {
				t[uctx.columnName] = ct
			} else {
				t = make(map[string]columnType)
				t[uctx.columnName] = ct
			}
			uctx.tables[uctx.filter.TableNameGet()] = t

			uctx.filter.ColumnAdd(uctx.columnName, traw)
		}
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

	d := append(uctx.security.tmpBuf, append(deferred, token...)...)

	uctx.security.isSkip = false
	uctx.security.tmpBuf = []byte{}

	// Check insert into table name
	if tn != uctx.filter.TableNameGet() {
		return d, fmt.Errorf("`create` and `insert into` table names are mismatch (create table: '%s', insert into table: '%s')", uctx.filter.TableNameGet(), tn)
	}

	return d, nil
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

	return rowDataGen(uctx), nil
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

	return rowDataGen(uctx), nil
}

func rowDataGen(uctx *userCtx) []byte {

	var out string

	row := uctx.filter.ValuePop()

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
