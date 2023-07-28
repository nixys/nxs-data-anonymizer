package relfilter

import (
	"bytes"
	"fmt"

	"github.com/nixys/nxs-data-anonymizer/misc"
)

type Rules struct {
	Tables map[string]TableRules
}

type TableRules struct {
	Columns map[string]ColumnRule
}

type ColumnRule struct {
	Value  string
	Unique bool
}

type Filter struct {

	// Rules for filter a table values
	rules Rules

	// Temp table data for filtering
	tableData tableData
}

type Row struct {
	Values []rowValue
}

type tableData struct {
	name    string
	columns map[string]int
	values  []rowValue
	uniques map[string]map[string]any
}

type rowValue struct {
	V []byte
	T ValueType
}

const uniqueAttempts = 5

func Init(rules Rules) *Filter {
	return &Filter{
		rules: rules,
	}
}

// TableCreate creates new data set for table `name`
func (filter *Filter) TableCreate(name string) {
	filter.tableData = tableData{
		name:    name,
		columns: make(map[string]int),
		uniques: make(map[string]map[string]any),
		values:  []rowValue{},
	}
}

func (filter *Filter) TableNameGet() string {
	return filter.tableData.name
}

// ColumnAdd adds new column into current data set
func (filter *Filter) ColumnAdd(name string) {
	filter.tableData.columns[name] = len(filter.tableData.columns)
}

func (filter *Filter) ColumnPop() {
	var lc string

	// Get last element index
	l := len(filter.tableData.columns) - 1
	if l < 0 {
		return
	}

	// Get last element column name
	for k, v := range filter.tableData.columns {
		if v == l {
			lc = k
		}
	}

	delete(filter.tableData.columns, lc)
}

func (filter *Filter) ValueByteAdd(b []byte) {
	filter.valueAdd(ValueTypeByte, b)
}

func (filter *Filter) ValueStringAdd(b []byte) {
	filter.valueAdd(ValueTypeString, b)
}

func (filter *Filter) ValueBinaryAdd(b []byte) {
	filter.valueAdd(ValueTypeBinary, b)
}

func (filter *Filter) ValueIntAdd(b []byte) {
	filter.valueAdd(ValueTypeInt, b)
}

func (filter *Filter) ValueFloatAdd(b []byte) {
	filter.valueAdd(ValueTypeFloat, b)
}

func (filter *Filter) ValueNULLAdd(b []byte) {
	filter.valueAdd(ValueTypeNULL, b)
}

// Apply applies filter rules for current data set
func (filter *Filter) Apply() error {

	type TplData struct {
		TableName string
		Values    map[string]string
	}

	tname := filter.tableData.name

	// Check current talbe exist in rules
	t, b := filter.rules.Tables[tname]
	if b == true {

		td := TplData{
			TableName: tname,
			Values:    make(map[string]string),
		}

		for c, n := range filter.tableData.columns {
			td.Values[c] = string(filter.tableData.values[n].V)
		}

		// Filter all columns with specified rules
		for cname, id := range filter.tableData.columns {

			// Check rule set for current column
			c, e := t.Columns[cname]
			if e == true {

				v, err := func() ([]byte, error) {

					for i := 0; i < uniqueAttempts; i++ {

						v, err := misc.TemplateExec(
							c.Value,
							td,
						)
						if err != nil {
							return []byte{}, fmt.Errorf("value compile template: %w", err)
						}

						v = bytes.ReplaceAll(v, []byte("\n"), []byte("\\n"))

						if c.Unique == false {
							return v, nil
						}

						var uv map[string]any
						if _, b := filter.tableData.uniques[cname]; b == false {
							// For first values
							uv = make(map[string]any)
						} else {
							uv = filter.tableData.uniques[cname]
						}

						if _, b := uv[string(v)]; b == false {
							uv[string(v)] = nil
							filter.tableData.uniques[cname] = uv
							return v, nil
						}
					}

					return []byte{}, fmt.Errorf("unable to generate unique value for column `%s.%s`, check filter value for this column in config", filter.tableData.name, cname)
				}()
				if err != nil {
					return err
				}

				// Set specified value in accordance with filter
				filter.tableData.values[id].V = v
			}
		}
	}

	return nil
}

// RowPull pulls row values for current data set.
// Row values will be dropped up after extract.
func (filter *Filter) RowPull() Row {

	// Save current values
	r := filter.tableData.values

	filter.rowDrop()

	return Row{
		Values: r,
	}
}

// rowDrop drops current row values
func (filter *Filter) rowDrop() {
	filter.tableData.values = []rowValue{}
}

func (filter *Filter) valueAdd(t ValueType, b []byte) {
	filter.tableData.values = append(
		filter.tableData.values,
		rowValue{
			V: bcopy(b),
			T: t,
		},
	)
}

// bcopy makes a bytes copy
func bcopy(b []byte) []byte {

	d := make([]byte, len(b))
	copy(d, b)

	return d
}
