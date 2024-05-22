package relfilter

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/nixys/nxs-data-anonymizer/misc"
)

type Rules struct {
	Tables           map[string]TableRules
	ExceptionColumns map[string]any
	Defaults         TableRules
	RandomizeTypes   map[ColumnType]ColumnRule
}

type TableRules struct {
	Columns map[string]ColumnRule
}

type ColumnRule struct {
	Type   misc.ValueType
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

const (
	envVarTable        = "ENVVARTABLE"
	envVarColumnPrefix = "ENVVARCOLUMN_"
	envVarCurColumn    = "ENVVARCURCOLUMN"
)

type rule struct {
	c  *column
	i  int
	cr ColumnRule
}

var RandomizeTypesDefault = map[ColumnType]ColumnRule{
	ColumnTypeBinary: {
		Type:   misc.ValueTypeTemplate,
		Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
		Unique: false,
	},
	ColumnTypeNum: {
		Type:   misc.ValueTypeTemplate,
		Value:  "0",
		Unique: false,
	},
	ColumnTypeString: {
		Type:   misc.ValueTypeTemplate,
		Value:  "randomized string data",
		Unique: false,
	},
}

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

// TableRulesLookup looks up filters for specified table name
func (filter *Filter) TableRulesLookup(name string) *TableRules {
	if t, b := filter.rules.Tables[name]; b {
		return &t
	}
	return nil
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

func (filter *Filter) Apply() error {

	var rls []rule

	tname := filter.tableData.name

	// Check rules exist for current table
	tr := filter.TableRulesLookup(tname)

	// Create rules for every column within current table
	for i, c := range filter.tableData.columns.cc {

		// Check direct rules for column
		if tr != nil {
			if cr, e := tr.Columns[c.n]; e == true {

				rls = append(
					rls,
					rule{
						c:  c,
						i:  i,
						cr: cr,
					},
				)
				continue
			}
		}

		// Check default rules for column
		if cr, e := filter.rules.Defaults.Columns[c.n]; e == true {
			rls = append(
				rls,
				rule{
					c:  c,
					i:  i,
					cr: cr,
				},
			)
			continue
		}

		// Check randomize rules for column
		if cr, b := filter.rules.RandomizeTypes[c.t]; b {

			// Check that column excepted
			if _, b := filter.rules.ExceptionColumns[c.n]; !b {
				rls = append(
					rls,
					rule{
						c:  c,
						i:  i,
						cr: cr,
					},
				)
				continue
			}
		}

		// Other rules if required
	}

	// Apply rules
	if err := filter.applyRules(tname, rls); err != nil {
		return fmt.Errorf("filters apply: %w", err)
	}

	return nil
}

func (filter *Filter) applyRules(tname string, rls []rule) error {

	// If no columns has rules
	if len(rls) == 0 {
		return nil
	}

	// Fill table data and table envs
	td := misc.TemplateData{
		TableName: tname,
		Values:    make(map[string][]byte),
	}

	tdenv := []string{
		fmt.Sprintf("%s=%s", envVarTable, tname),
	}

	for i, c := range filter.tableData.columns.cc {
		td.Values[c.n] = filter.tableData.values[i].V

		tdenv = append(
			tdenv,
			fmt.Sprintf("%s%s=%s", envVarColumnPrefix, c.n, string(filter.tableData.values[i].V)),
		)
	}

	// Apply rule for each specified column
	for _, r := range rls {

		var tde []string

		// Create tmp env variables with current column name
		tde = append(
			tdenv,
			fmt.Sprintf("%s=%s", envVarCurColumn, r.c.n),
		)

		v, err := filter.applyFilter(r.c.n, r.cr, td, tde)
		if err != nil {
			return fmt.Errorf("rules: %w", err)
		}

		// Set specified value in accordance with filter
		filter.tableData.values[r.i].V = v
	}

	return nil
}

func (filter *Filter) applyFilter(cn string, cr ColumnRule, td misc.TemplateData, tde []string) ([]byte, error) {

	for i := 0; i < uniqueAttempts; i++ {

		var (
			v   []byte
			err error
		)

		switch cr.Type {
		case misc.ValueTypeTemplate:
			v, err = misc.TemplateExec(
				cr.Value,
				td,
			)
			if err != nil {
				return []byte{}, fmt.Errorf("filter: value compile template: %w", err)
			}
		case misc.ValueTypeCommand:

			var stderr, stdout bytes.Buffer

			cmd := exec.Command(cr.Value)

			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			cmd.Env = tde

			if err := cmd.Run(); err != nil {

				e, b := err.(*exec.ExitError)
				if b == false {
					return []byte{}, fmt.Errorf("filter: value exec command: %w", err)
				}

				return []byte{}, fmt.Errorf("filter: value exec command: bad exit code %d: %s", e.ExitCode(), stderr.String())
			}

			v = stdout.Bytes()

		default:
			return []byte{}, fmt.Errorf("filter: value compile: unknown type")
		}

		v = bytes.ReplaceAll(v, []byte("\n"), []byte("\\n"))

		if cr.Unique == false {
			return v, nil
		}

		var uv map[string]any
		if _, b := filter.tableData.uniques[cn]; b == false {
			// For first values
			uv = make(map[string]any)
		} else {
			uv = filter.tableData.uniques[cn]
		}

		if _, b := uv[string(v)]; b == false {
			uv[string(v)] = nil
			filter.tableData.uniques[cn] = uv
			return v, nil
		}
	}

	return []byte{}, fmt.Errorf("filter: unable to generate unique value for column `%s.%s`, check filter value for this column in config", filter.tableData.name, cn)
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
