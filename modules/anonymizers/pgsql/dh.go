package pgsql_anonymize

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func dhTableName(usrCtx any, deferred, token []byte) ([]byte, error) {

	tname := bytes.TrimSpace(deferred)

	filter := usrCtx.(*relfilter.Filter)
	filter.TableCreate(string(tname))

	return append(deferred, token...), nil
}

func dhFieldName(usrCtx any, deferred, token []byte) ([]byte, error) {

	fname := bytes.TrimSpace(deferred)

	filter := usrCtx.(*relfilter.Filter)
	filter.ColumnAdd(string(fname))

	return append(deferred, token...), nil
}

func dhValue(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	valuePut(filter, deferred, token)

	return []byte{}, nil
}

func dhValueEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	valuePut(filter, deferred, token)

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
			out += "\t"
		}

		switch v.T {
		default:
			out += fmt.Sprintf("%s", v.V)
		}
	}

	return []byte(fmt.Sprintf("%s\n", out))
}

func valuePut(filter *relfilter.Filter, deferred, token []byte) {
	if _, b := strconv.ParseInt(string(deferred), 10, 64); b == nil {
		filter.ValueIntAdd(deferred)
	} else if _, b := strconv.ParseFloat(string(deferred), 10); b == nil {
		filter.ValueFloatAdd(deferred)
	} else {
		filter.ValueStringAdd(deferred)
	}
}
