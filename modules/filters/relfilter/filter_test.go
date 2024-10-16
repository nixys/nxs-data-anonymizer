package relfilter

import (
	"bytes"
	"testing"

	"github.com/nixys/nxs-data-anonymizer/misc"
)

func TestExecFilter(t *testing.T) {

	// Test `drop` template function
	v, d, err := execFilter(
		execFilterOpts{
			t: misc.ValueTypeTemplate,
			v: "{{- drop -}}",
		},
		nil,
		nil)
	if err != nil {
		t.Fatal("`drop` function:", err)
	}
	if v != nil || d == false {
		t.Fatal("`drop` function: incorrect return value")
	}
	t.Logf("`drop` function: success")

	// Test `null` template function
	v, d, err = execFilter(
		execFilterOpts{
			t: misc.ValueTypeTemplate,
			v: "{{- null -}}",
		},
		nil,
		nil)
	if err != nil {
		t.Fatal("`null` function:", err)
	}
	if v != nil || d == true {
		t.Fatal("`null` function: incorrect return value")
	}
	t.Logf("`null` function: success")
}

func TestFilterApply(t *testing.T) {
	TestFilterApplyDropFunction(t)
	TestFilterApplyNullFunction(t)
	TestLinkFilterApplyDropFunction(t)
}

func TestFilterApplyDropFunction(t *testing.T) {

	f, err := Init(
		InitOpts{
			TableRules: map[string]map[string]ColumnRuleOpts{
				"testTable": {
					"testColumn1": ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- drop -}}",
						Unique: false,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init:", err)
	}

	testFilterTableInit(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r := f.ValuePop()

	if r.Values != nil {
		t.Fatal("`drop` function: unexpected behaviour")
	}
	t.Logf("`drop` function: success")
}

func TestFilterApplyNullFunction(t *testing.T) {

	f, err := Init(
		InitOpts{
			TableRules: map[string]map[string]ColumnRuleOpts{
				"testTable": {
					"testColumn1": ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- null -}}",
						Unique: false,
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init:", err)
	}

	testFilterTableInit(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r := f.ValuePop()

	if len(r.Values) < 2 || r.Values[0].V != nil {
		t.Fatal("`null` function: unexpected behaviour")
	}
	t.Logf("`null` function: success")
}

func TestLinkFilterApply(t *testing.T) {

	f, err := Init(
		InitOpts{
			TableRules: map[string]map[string]ColumnRuleOpts{
				"testTable1": {
					"testColumn1": ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- 11 -}}",
						Unique: false,
					},
				},
				"testTable2": {
					"testColumn1": ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- 22 -}}",
						Unique: false,
					},
				},
			},

			Link: []LinkOpts{
				{
					Rule: ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- uuidv4 -}}",
						Unique: false,
					},
					With: map[string][]string{
						"testTable1": {
							"testColumn2",
						},
						"testTable2": {
							"testColumn2",
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init:", err)
	}

	// Fill table 1
	testLinkFilterTable1Init(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r1 := f.ValuePop()

	if len(r1.Values) < 2 {
		t.Fatal("incorrect row len for table 1")
	}

	// Fill table 2
	testLinkFilterTable2Init(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r2 := f.ValuePop()

	if len(r2.Values) < 2 {
		t.Fatal("incorrect row len for table 2")
	}

	if bytes.Compare(r1.Values[1].V, r2.Values[1].V) != 0 {
		t.Fatal("incorrect values for tables after filter apply")
	}

	t.Logf("success")
}

func TestLinkFilterApplyDropFunction(t *testing.T) {

	f, err := Init(
		InitOpts{
			TableRules: map[string]map[string]ColumnRuleOpts{
				"testTable1": {
					"testColumn1": ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- 1 -}}",
						Unique: false,
					},
				},
			},

			Link: []LinkOpts{
				{
					Rule: ColumnRuleOpts{
						Type:   misc.ValueTypeTemplate,
						Value:  "{{- drop -}}",
						Unique: false,
					},
					With: map[string][]string{
						"testTable1": {
							"testColumn2",
						},
						"testTable2": {
							"testColumn2",
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init:", err)
	}

	// Fill table 1
	testLinkFilterTable1Init(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r1 := f.ValuePop()

	if r1.Values != nil {
		t.Fatal("`drop` function: unexpected behaviour for table 1")
	}

	// Fill table 2
	testLinkFilterTable2Init(f)

	// Apply filters for row
	if err := f.Apply(); err != nil {
		t.Fatal("apply:", err)
	}

	// Get row values
	r2 := f.ValuePop()

	if r2.Values != nil {
		t.Fatal("`drop` function: unexpected behaviour for table 2")
	}

	t.Logf("`drop` function: success")
}

func testFilterTableInit(f *Filter) {
	// Create table
	f.TableCreate("testTable")

	// Add column `testColumn1` with value
	f.ColumnAdd("testColumn1", "int")
	f.ValueAdd([]byte("0"))

	// Add column `testColumn2` with value
	f.ColumnAdd("testColumn2", "varchar(10)")
	f.ValueAdd([]byte("a"))
}

func testLinkFilterTable1Init(f *Filter) {
	// Create table
	f.TableCreate("testTable1")

	// Add column `testColumn1` with value
	f.ColumnAdd("testColumn1", "int")
	f.ValueAdd([]byte("1"))

	// Add column `testColumn2` with value
	f.ColumnAdd("testColumn2", "varchar(100)")
	f.ValueAdd([]byte("a"))
}

func testLinkFilterTable2Init(f *Filter) {
	// Create table
	f.TableCreate("testTable2")

	// Add column `testColumn1` with value
	f.ColumnAdd("testColumn1", "int")
	f.ValueAdd([]byte("2"))

	// Add column `testColumn2` with value
	f.ColumnAdd("testColumn2", "varchar(100)")
	f.ValueAdd([]byte("a"))
}
