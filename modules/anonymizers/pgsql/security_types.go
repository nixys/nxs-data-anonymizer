package pgsql_anonymize

import (
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

const (
	securityTypeInt    = "0"
	securityTypeFloat  = "0.0"
	securityTypeString = "randomized string data"
)

var typeRuleDefault = []relfilter.TypeRuleOpts{

	// Integer
	{
		Selector: "(?i)^smallint",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^integer",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bigint",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^smallserial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^serial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bigserial",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeInt,
			Unique: false,
		},
	},

	// Float
	{
		Selector: "(?i)^decimal",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeFloat,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^numeric",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeFloat,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^real",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeFloat,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^double",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeFloat,
			Unique: false,
		},
	},

	// Strings
	{
		Selector: "(?i)^character",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeString,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bpchar",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeString,
			Unique: false,
		},
	},
	{
		Selector: "(?i)^text",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  securityTypeString,
			Unique: false,
		},
	},
}
