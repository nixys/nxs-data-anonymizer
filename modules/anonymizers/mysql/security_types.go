package mysql_anonymize

import (
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
)

var typeRuleDefault = []relfilter.TypeRuleOpts{

	// String
	{
		Selector: "(?i)^char\\((\\d+)\\)|^char",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{{ trunc (index .ColumnTypeGroups 0 1 | int ) \"randomized char\" }}",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^varchar\\((\\d+)\\)|^varchar",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{{ trunc (index .ColumnTypeGroups 0 1 | int ) \"randomized varchar\" }}",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^tinytext",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized tinytext",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^text",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized text",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^mediumtext",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized mediumtext",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^longtext",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "randomized longtext",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^enum\\(.*'(.*)'.*\\)",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{{ index .ColumnTypeGroups 0 1 }}",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^set\\(.*'(.*)'.*\\)",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{{ index .ColumnTypeGroups 0 1 }}",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^datetime",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "2024-01-01 00:00:00",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^date",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "2024-01-01",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^timestamp",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "2024-01-01 00:00:00",
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
	{
		Selector: "(?i)^year",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "2024",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^json",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "{\"randomized\": \"json_data\"}",
			Unique: false,
		},
	},

	// Numeric
	{
		Selector: "(?i)^bit",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^bool",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^boolean",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^tinyint",
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
		Selector: "(?i)^mediumint",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^int",
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
		Selector: "(?i)^float",
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
	{
		Selector: "(?i)^decimal",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^dec",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "0.0",
			Unique: false,
		},
	},

	// Binary
	{
		Selector: "(?i)^binary",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^varbinary",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	}, {
		Selector: "(?i)^tinyblob",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^blob",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	},
	{
		Selector: "(?i)^mediumblob",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	}, {
		Selector: "(?i)^longblob",
		Rule: relfilter.ColumnRuleOpts{
			Type:   misc.ValueTypeTemplate,
			Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
			Unique: false,
		},
	},
}
