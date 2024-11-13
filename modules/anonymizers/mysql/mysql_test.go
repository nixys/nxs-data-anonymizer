package mysql_anonymize

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func TestMySQL(t *testing.T) {

	var r, e bytes.Buffer

	fin, err := os.Open(".testdata/mysql_test.in.sql")
	if err != nil {
		t.Fatal("open input SQL:", err)
	}

	m, err := Init(
		fin,
		InitOpts{
			Rules: RulesOpts{
				TableRules: map[string]map[string]relfilter.ColumnRuleOpts{

					// Delete only row with id `2`
					"table1": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.id \"2\" }}{{ drop }}{{ else }}{{ .Values.id }}{{ end }}",
							Unique: false,
						},
					},

					// Delete all rows from table
					"table2": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ drop }}",
							Unique: false,
						},
					},

					// Delete no rows
					"table3": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.id \"4\" }}{{ drop }}{{ else }}{{ .Values.id }}{{ end }}",
							Unique: false,
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init MySQL:", err)
	}

	if err := m.Run(context.Background(), &r); err != nil {
		t.Fatal("run MySQL:", err)
	}

	fout, err := os.Open(".testdata/mysql_test.out.sql")
	if err != nil {
		t.Fatal("open output SQL:", err)
	}

	if _, err := e.ReadFrom(fout); err != nil {
		t.Fatal("read output SQL:", err)
	}

	// os.WriteFile(".testdata/mysql_test.out.sql", r.Bytes(), 0644)

	if r.String() != e.String() {
		t.Fatal("incorrect anonymization result")
	}

	t.Logf("success")
}

func TestMySQLDos(t *testing.T) {

	var r, e bytes.Buffer

	fin, err := os.Open(".testdata/mysql_test.dos.in.sql")
	if err != nil {
		t.Fatal("open input SQL:", err)
	}

	m, err := Init(
		fin,
		InitOpts{
			Rules: RulesOpts{
				TableRules: map[string]map[string]relfilter.ColumnRuleOpts{

					// Delete only row with id `2`
					"table1": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.id \"2\" }}{{ drop }}{{ else }}{{ .Values.id }}{{ end }}",
							Unique: false,
						},
					},

					// Delete all rows from table
					"table2": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ drop }}",
							Unique: false,
						},
					},

					// Delete no rows
					"table3": {
						"id": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.id \"4\" }}{{ drop }}{{ else }}{{ .Values.id }}{{ end }}",
							Unique: false,
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init MySQL:", err)
	}

	if err := m.Run(context.Background(), &r); err != nil {
		t.Fatal("run MySQL:", err)
	}

	fout, err := os.Open(".testdata/mysql_test.dos.out.sql")
	if err != nil {
		t.Fatal("open output SQL:", err)
	}

	if _, err := e.ReadFrom(fout); err != nil {
		t.Fatal("read output SQL:", err)
	}

	// os.WriteFile(".testdata/mysql_test.dos.out.sql", r.Bytes(), 0644)

	if r.String() != e.String() {
		t.Fatal("incorrect anonymization result")
	}

	t.Logf("success")
}
