package mysql_anonymize

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhCreateTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	filter.TableCreate(string(deferred))

	return append(deferred, token...), nil
}

func dhCreateTableFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	filter.ColumnAdd(string(deferred))

	return append(deferred, token...), nil
}

func dhInsertIntoTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)

	// Check insert into table name
	if bytes.Compare([]byte(filter.TableNameGet()), deferred) != 0 {
		return append(deferred, token...), fmt.Errorf("`create` and `insert into` table names are mismatch (create table: '%s', insert into table: '%s')", filter.TableNameGet(), string(deferred))
	}

	return append(deferred, token...), nil
}

func dhCreateTableValues(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	valuePut(filter, deferred, token)

	return []byte{}, nil
}

func dhCreateTableValuesBinary(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	filter.ValueBinaryAdd(deferred)

	return []byte{}, nil
}

func dhCreateTableValuesEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	valuePut(filter, deferred, token)

	// Apply filter for row
	if err := filter.Apply(); err != nil {
		return []byte{}, err
	}

	return rowDataGen(filter), nil
}

func dhCreateTableValuesStringEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)

	// Apply filter for row
	if err := filter.Apply(); err != nil {
		return []byte{}, err
	}

	return rowDataGen(filter), nil
}

func rowDataGen(filter *relfilter.Filter) []byte {

	var out string

	row := filter.RowPull()

	for i, v := range row.Values {

		if i > 0 {
			out += ","
		}

		switch v.T {
		case relfilter.ValueTypeString:
			out += fmt.Sprintf("'%s'", v.V)
		case relfilter.ValueTypeBinary:
			out += fmt.Sprintf("_binary '%s'", v.V)
		default:
			out += fmt.Sprintf("%s", v.V)
		}
	}

	return []byte(fmt.Sprintf("(%s)", out))
}

func valuePut(filter *relfilter.Filter, deferred, token []byte) {
	if bytes.Compare(token, []byte{'\''}) == 0 {
		filter.ValueStringAdd(deferred)
	} else if bytes.Compare(deferred, []byte("NULL")) == 0 {
		filter.ValueNULLAdd(deferred)
	} else if _, b := strconv.ParseInt(string(deferred), 10, 64); b == nil {
		filter.ValueIntAdd(deferred)
	} else if _, b := strconv.ParseFloat(string(deferred), 10); b == nil {
		filter.ValueFloatAdd(deferred)
	} else {
		filter.ValueByteAdd(deferred)
	}
}
