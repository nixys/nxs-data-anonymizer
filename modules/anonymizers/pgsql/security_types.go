package pgsql_anonymize

import (
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

var typeRuleDefault = []relfilter.TypeRuleOpts{

	// Integer
	{
		Selector: "(?i)^boolean",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^smallint",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^integer",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bigint",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^smallserial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^serial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bigserial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},

	// Float
	{
		Selector: "(?i)^decimal",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^numeric",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^real",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^double",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},

	// Strings
	{
		Selector: "(?i)^character",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized character data",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bpchar",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized bpchar data",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^text",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized text data",
			Unique: false,
		},
	},

	// Date & time
	{
		Selector: "(?i)^date",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "2024-01-01",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^time",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "00:00:00",
			Unique: false,
		},
	},

	// Structures
	{
		Selector: "(?i)^jsonb",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{\"randomized\": \"json_data\"}",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^xml",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "<note><body>Randomized XML</body></note>",
			Unique: false,
		},
	},
}
