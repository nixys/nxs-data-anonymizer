package relfilter

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/nixys/nxs-data-anonymizer/misc"
)

type InitOpts struct {
	Variables map[string]VariableRuleOpts

	TableRules       map[string]map[string]ColumnRuleOpts
	DefaultRules     map[string]ColumnRuleOpts
	ExceptionColumns []string

	TypeRuleCustom  []TypeRuleOpts
	TypeRuleDefault []TypeRuleOpts
}

type TypeRuleOpts struct {
	Selector string
	Rule     ColumnRuleOpts
}

type ColumnRuleOpts struct {
	Type   misc.ValueType
	Value  string
	Unique bool
}

type VariableRuleOpts struct {
	Type  misc.ValueType
	Value string
}

type Filter struct {

	// Rules for filter a table values
	rules rules

	// Temp table data for filtering
	tableData tableData
}

type Row struct {
	Values []rowValue
}

type rules struct {
	variables map[string]string

	tableRules       map[string]map[string]ColumnRuleOpts
	defaultRules     map[string]ColumnRuleOpts
	exceptionColumns map[string]any

	typeRuleCustom  []typeRule
	typeRuleDefault []typeRule
}

type typeRule struct {
	Rgx  *regexp.Regexp
	Rule ColumnRuleOpts
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
	envVarGlobalPrefix = "ENVVARGLOBAL_"
	envVarTable        = "ENVVARTABLE"
	envVarColumnPrefix = "ENVVARCOLUMN_"
	envVarCurColumn    = "ENVVARCURCOLUMN"
)

type applyRule struct {
	c  *column
	i  int
	cr ColumnRuleOpts
}

func Init(opts InitOpts) (*Filter, error) {

	trc := []typeRule{}
	trd := []typeRule{}

	// Make custom type rules
	for _, r := range opts.TypeRuleCustom {

		re, err := regexp.Compile(r.Selector)
		if err != nil {
			return nil, fmt.Errorf("filter init: %w", err)
		}

		trc = append(
			trc,
			typeRule{
				Rgx:  re,
				Rule: r.Rule,
			},
		)
	}

	// Make default type rules
	for _, r := range opts.TypeRuleDefault {

		re, err := regexp.Compile(r.Selector)
		if err != nil {
			return nil, fmt.Errorf("filter init: %w", err)
		}

		trd = append(
			trd,
			typeRule{
				Rgx:  re,
				Rule: r.Rule,
			},
		)
	}

	// Make exceptions
	excpts := make(map[string]any)
	for _, e := range opts.ExceptionColumns {
		excpts[e] = nil
	}

	vars := make(map[string]string)
	for n, f := range opts.Variables {
		v, err := makeVariable(f)
		if err != nil {
			return nil, fmt.Errorf("filter init: %w", err)
		}
		vars[n] = v
	}

	return &Filter{
		rules: rules{
			variables:        vars,
			tableRules:       opts.TableRules,
			defaultRules:     opts.DefaultRules,
			exceptionColumns: excpts,
			typeRuleCustom:   trc,
			typeRuleDefault:  trd,
		},
	}, nil
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
func (filter *Filter) TableRulesLookup(name string) map[string]ColumnRuleOpts {
	if t, b := filter.rules.tableRules[name]; b {
		return t
	}
	return nil
}

// ColumnAdd adds new column into current data set
func (filter *Filter) ColumnAdd(name string, rt string) {
	filter.tableData.columns.add(name, rt)
}

func (filter *Filter) ColumnGetName(index int) string {
	return filter.tableData.columns.getNameByIndex(index)
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

	var rls []applyRule

	tname := filter.tableData.name

	// Check rules exist for current table
	tr := filter.TableRulesLookup(tname)

	// Create rules for every column within current table
	for i, c := range filter.tableData.columns.cc {

		// Check direct rules for column
		if tr != nil {
			if cr, e := tr[c.n]; e == true {

				rls = append(
					rls,
					applyRule{
						c:  c,
						i:  i,
						cr: cr,
					},
				)
				continue
			}
		}

		// Check default rules for column
		if cr, e := filter.rules.defaultRules[c.n]; e == true {
			rls = append(
				rls,
				applyRule{
					c:  c,
					i:  i,
					cr: cr,
				},
			)
			continue
		}

		// Check column is excepted
		if _, b := filter.rules.exceptionColumns[c.n]; b {
			continue
		}

		// Check custom type rule for column
		if b := func() bool {
			for _, r := range filter.rules.typeRuleCustom {
				if r.Rgx.Match([]byte(c.rawType)) {
					rls = append(
						rls,
						applyRule{
							c:  c,
							i:  i,
							cr: r.Rule,
						},
					)
					return true
				}
			}
			return false
		}(); b {
			continue
		}

		// Check default type rule for column
		if b := func() bool {
			for _, r := range filter.rules.typeRuleDefault {
				if r.Rgx.Match([]byte(c.rawType)) {
					rls = append(
						rls,
						applyRule{
							c:  c,
							i:  i,
							cr: r.Rule,
						},
					)
					return true
				}
			}
			return false
		}(); b {
			continue
		}

		// Other rules if required
	}

	// Apply rules
	if err := filter.applyRules(tname, rls); err != nil {
		return fmt.Errorf("filters apply: %w", err)
	}

	return nil
}

func (filter *Filter) applyRules(tname string, rls []applyRule) error {

	// If no columns has rules
	if len(rls) == 0 {
		return nil
	}

	// Fill table data and table envs
	td := misc.TemplateData{
		TableName: tname,
		Values:    make(map[string][]byte),
		Variables: filter.rules.variables,
	}

	tdenv := []string{
		fmt.Sprintf("%s=%s", envVarTable, tname),
	}

	for n, v := range filter.rules.variables {
		tdenv = append(
			tdenv,
			fmt.Sprintf("%s%s=%s", envVarGlobalPrefix, n, v),
		)
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

		v, err := filter.applyColumnFilter(r.c.n, r.cr, td, tde)
		if err != nil {
			return fmt.Errorf("rules: %w", err)
		}

		// Set specified value in accordance with filter
		filter.tableData.values[r.i].V = v
	}

	return nil
}

func (filter *Filter) applyColumnFilter(cn string, cr ColumnRuleOpts, td misc.TemplateData, tde []string) ([]byte, error) {

	for i := 0; i < uniqueAttempts; i++ {

		var (
			v   []byte
			err error
		)

		switch cr.Type {
		case misc.ValueTypeTemplate:
			v, err = misc.TemplateExec(
				cr.Value,
				&td,
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

func makeVariable(cr VariableRuleOpts) (string, error) {

	var (
		v   []byte
		err error
	)

	switch cr.Type {
	case misc.ValueTypeTemplate:
		v, err = misc.TemplateExec(
			cr.Value,
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("variable: value compile template: %w", err)
		}
	case misc.ValueTypeCommand:

		var stderr, stdout bytes.Buffer

		cmd := exec.Command(cr.Value)

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {

			e, b := err.(*exec.ExitError)
			if b == false {
				return "", fmt.Errorf("variable: value exec command: %w", err)
			}

			return "", fmt.Errorf("variable: value exec command: bad exit code %d: %s", e.ExitCode(), stderr.String())
		}

		v = stdout.Bytes()

	default:
		return "", fmt.Errorf("variable: value compile: unknown type")
	}

	return string(bytes.ReplaceAll(v, []byte("\n"), []byte("\\n"))), nil
}
