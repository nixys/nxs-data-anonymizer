package pgsql_anonymize

import (
	"bytes"
	"fmt"

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
	filter.ColumnAdd(string(fname), relfilter.ColumnTypeNone)

	return append(deferred, token...), nil
}

func dhValue(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	filter.ValueAdd(deferred)

	return []byte{}, nil
}

func dhValueEnd(usrCtx any, deferred, token []byte) ([]byte, error) {

	filter := usrCtx.(*relfilter.Filter)
	filter.ValueAdd(deferred)

	// Apply filter for row
	if err := filter.Apply(); err != nil {
		return []byte{}, err
	}

	return rowDataGen(filter), nil
}

func rowDataGen(filter *relfilter.Filter) []byte {

	var out string

	row := filter.ValuePop()

	for i, v := range row.Values {

		if i > 0 {
			out += "\t"
		}

		out += fmt.Sprintf("%s", v.V)
	}

	return []byte(fmt.Sprintf("%s\n", out))
}
