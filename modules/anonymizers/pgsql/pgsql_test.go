package pgsql_anonymize

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

func TestPgSQL(t *testing.T) {

	var r, e bytes.Buffer

	fin, err := os.Open(".testdata/pgsql_test.in.sql")
	if err != nil {
		t.Fatal("open input SQL:", err)
	}

	m, err := Init(
		fin,
		InitOpts{
			Rules: RulesOpts{
				TableRules: map[string]map[string]relfilter.ColumnRuleOpts{

					// Delete only row with id `2`
					"public.list_types": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.integer_type \"8765542\" }}{{ drop }}{{ else }}{{ .Values.integer_type }}{{ end }}",
							Unique: false,
						},
					},

					// Delete all rows from table
					"public.list_types2": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ drop }}",
							Unique: false,
						},
					},

					// Delete no rows
					"public.list_types3": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.integer_type \"0\" }}{{ drop }}{{ else }}{{ .Values.integer_type }}{{ end }}",
							Unique: false,
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init PgSQL:", err)
	}

	if err := m.Run(context.Background(), &r); err != nil {
		t.Fatal("run PgSQL:", err)
	}

	fout, err := os.Open(".testdata/pgsql_test.out.sql")
	if err != nil {
		t.Fatal("open output SQL:", err)
	}

	if _, err := e.ReadFrom(fout); err != nil {
		t.Fatal("read output SQL:", err)
	}

	// os.WriteFile(".testdata/pgsql_test.out.sql", r.Bytes(), 0644)

	if r.String() != e.String() {
		t.Fatal("incorrect anonymization result")
	}

	t.Logf("success")
}

func TestPgSQLDos(t *testing.T) {

	var r, e bytes.Buffer

	fin, err := os.Open(".testdata/pgsql_test.dos.in.sql")
	if err != nil {
		t.Fatal("open input SQL:", err)
	}

	m, err := Init(
		fin,
		InitOpts{
			Rules: RulesOpts{
				TableRules: map[string]map[string]relfilter.ColumnRuleOpts{

					// Delete only row with id `2`
					"public.list_types": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.integer_type \"8765542\" }}{{ drop }}{{ else }}{{ .Values.integer_type }}{{ end }}",
							Unique: false,
						},
					},

					// Delete all rows from table
					"public.list_types2": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ drop }}",
							Unique: false,
						},
					},

					// Delete no rows
					"public.list_types3": {
						"integer_type": relfilter.ColumnRuleOpts{
							Type:   misc.ValueTypeTemplate,
							Value:  "{{ if eq .Values.integer_type \"0\" }}{{ drop }}{{ else }}{{ .Values.integer_type }}{{ end }}",
							Unique: false,
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal("init PgSQL:", err)
	}

	if err := m.Run(context.Background(), &r); err != nil {
		t.Fatal("run PgSQL:", err)
	}

	fout, err := os.Open(".testdata/pgsql_test.dos.out.sql")
	if err != nil {
		t.Fatal("open output SQL:", err)
	}

	if _, err := e.ReadFrom(fout); err != nil {
		t.Fatal("read output SQL:", err)
	}

	// os.WriteFile(".testdata/pgsql_test.dos.out.sql", r.Bytes(), 0644)

	if r.String() != e.String() {
		t.Fatal("incorrect anonymization result")
	}

	t.Logf("success")
}
