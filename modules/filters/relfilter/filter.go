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
	columns columns
	values  []rowValue
	uniques map[string]map[string]any
}

type rowValue struct {
	V []byte
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
		columns: columnsInit(),
		uniques: make(map[string]map[string]any),
		values:  []rowValue{},
	}
}

func (filter *Filter) TableNameGet() string {
	return filter.tableData.name
}

// ColumnAdd adds new column into current data set
func (filter *Filter) ColumnAdd(name string, t ColumnType) {
	filter.tableData.columns.add(name, t)
}

func (filter *Filter) ColumnTypeGet(index int) ColumnType {
	return filter.tableData.columns.typeGetByIndex(index)
}

func (filter *Filter) ValueAdd(b []byte) {
	filter.tableData.values = append(
		filter.tableData.values,
		rowValue{
			V: bcopy(b),
		},
	)
}

// ValuePop pops the last values row from current data set
func (filter *Filter) ValuePop() Row {

	// Save current values
	r := filter.tableData.values

	filter.rowCleanup()

	return Row{
		Values: r,
	}
}

// Apply applies filter rules for current data set
func (filter *Filter) Apply() error {

	tname := filter.tableData.name

	// Check current table exist in rules
	t, b := filter.rules.Tables[tname]
	if b == true {

		td := misc.TemplateData{
			TableName: tname,
			Values:    make(map[string][]byte),
		}

		for i, d := range filter.tableData.columns.cc {
			td.Values[d.n] = filter.tableData.values[i].V
		}

		// Filter all columns with specified rules
		for n, d := range filter.tableData.columns.cc {

			// Check rule set for current column
			c, e := t.Columns[d.n]
			if e == false {
				continue
			}

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
					if _, b := filter.tableData.uniques[d.n]; b == false {
						// For first values
						uv = make(map[string]any)
					} else {
						uv = filter.tableData.uniques[d.n]
					}

					if _, b := uv[string(v)]; b == false {
						uv[string(v)] = nil
						filter.tableData.uniques[d.n] = uv
						return v, nil
					}
				}

				return []byte{}, fmt.Errorf("unable to generate unique value for column `%s.%s`, check filter value for this column in config", filter.tableData.name, d.n)
			}()
			if err != nil {
				return err
			}

			// Set specified value in accordance with filter
			filter.tableData.values[n].V = v
		}
	}

	return nil
}

// rowCleanup cleanups current row values
func (filter *Filter) rowCleanup() {
	filter.tableData.values = []rowValue{}
}

// bcopy makes a bytes copy
func bcopy(b []byte) []byte {

	if b == nil {
		return nil
	}

	d := make([]byte, len(b))
	copy(d, b)

	return d
}
