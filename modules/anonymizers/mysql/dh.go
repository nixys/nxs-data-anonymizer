package mysql_anonymize

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhCreateTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)
	uctx.filter.TableCreate(string(deferred))

	return append(deferred, token...), nil
}

func dhCreateTableFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)
	uctx.column.name = string(deferred)

	return append(deferred, token...), nil
}

func dhCreateTableColumnTypeAdd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

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

	if uctx.column.isSkip == false {
		uctx.filter.ColumnAdd(uctx.column.name, uctx.column.columnType)
	}

	uctx.column = userColumnCtx{}

	return append(deferred, token...), nil
}

func dhInsertIntoTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	// Check insert into table name
	if bytes.Compare([]byte(uctx.filter.TableNameGet()), deferred) != 0 {
		return append(deferred, token...), fmt.Errorf("`create` and `insert into` table names are mismatch (create table: '%s', insert into table: '%s')", uctx.filter.TableNameGet(), string(deferred))
	}

	return append(deferred, token...), nil
}

func dhCreateTableValues(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	if bytes.Compare(deferred, []byte("NULL")) == 0 {
		uctx.filter.ValueAdd(nil)
	} else {
		uctx.filter.ValueAdd(deferred)
	}

	return []byte{}, nil
}

func dhCreateTableValuesString(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

	uctx.filter.ValueAdd(deferred)

	return []byte{}, nil
}

func dhCreateTableValuesEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	uctx := usrCtx.(*userCtx)

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
