package pgsql_anonymize

import (
	"context"
	"io"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"

	fsm "github.com/nixys/nxs-go-fsm"
)

type InitOpts struct {
	Security SecurityOpts
	Rules    relfilter.Rules
}

type SecurityOpts struct {
	TablePolicy     misc.SecurityPolicyTablesType
	TableExceptions map[string]any
}

type userCtx struct {
	filter *relfilter.Filter

	security securityCtx

	tn     *string
	tables map[string]map[string]relfilter.ColumnType
}

type securityCtx struct {
	tmpBuf []byte
	isSkip bool

	tablePolicy     misc.SecurityPolicyTablesType
	tableExceptions map[string]any
}

const (
	columnTypeString relfilter.ColumnType = "string"
	columnTypeInt    relfilter.ColumnType = "integer"
	columnTypeFloat  relfilter.ColumnType = "float"
)

var typeKeys = map[string]relfilter.ColumnType{

	// Integer
	"smallint":    columnTypeInt,
	"integer":     columnTypeInt,
	"bigint":      columnTypeInt,
	"smallserial": columnTypeInt,
	"serial":      columnTypeInt,
	"bigserial":   columnTypeInt,

	// Float
	"decimal": columnTypeFloat,
	"numeric": columnTypeFloat,
	"real":    columnTypeFloat,
	"double":  columnTypeFloat,

	// Strings
	"character": columnTypeString,
	"bpchar":    columnTypeString,
	"text":      columnTypeString,
}

var RandomizeTypesDefault = map[relfilter.ColumnType]relfilter.ColumnRule{
	columnTypeInt: {
		Type:   misc.ValueTypeTemplate,
		Value:  "0",
		Unique: false,
	},
	columnTypeFloat: {
		Type:   misc.ValueTypeTemplate,
		Value:  "0.0",
		Unique: false,
	},
	columnTypeString: {
		Type:   misc.ValueTypeTemplate,
		Value:  "randomized string data",
		Unique: false,
	},
}

func userCtxInit(s InitOpts) *userCtx {
	return &userCtx{
		filter: relfilter.Init(s.Rules),
		tables: make(map[string]map[string]relfilter.ColumnType),
		security: securityCtx{
			tablePolicy:     s.Security.TablePolicy,
			tableExceptions: s.Security.TableExceptions,
		},
	}
}

func Init(ctx context.Context, r io.Reader, s InitOpts) io.Reader {

	return fsm.Init(
		r,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   userCtxInit(s),
			InitState: stateInit,
			States: map[fsm.StateName]fsm.State{

				stateInit: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableName,
							Switch: fsm.Switch{
								Trigger: []byte("COPY"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: dhSecurityCopy,
						},
						{
							Name: stateCreateTableName,
							Switch: fsm.Switch{
								Trigger: []byte("CREATE TABLE"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
					},
				},

				stateCreateTableName: {
					NextStates: []fsm.NextState{
						{
							Name: stateCreateTableTail,
							Switch: fsm.Switch{
								Trigger: []byte("("),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableName,
						},
					},
				},

				stateCreateTableTail: {
					NextStates: []fsm.NextState{
						{
							Name: stateInit,
							Switch: fsm.Switch{
								Trigger: []byte(");"),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableDesc,
						},
					},
				},

				stateTableName: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldName,
							Switch: fsm.Switch{
								Trigger: []byte("("),
							},
							DataHandler: dhTableName,
						},
					},
				},
				stateFieldName: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldName,
							Switch: fsm.Switch{
								Trigger: []byte(","),
							},
							DataHandler: dhFieldName,
						},
						{
							Name: stateCopyTail,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
							},
							DataHandler: dhFieldName,
						},
					},
				},
				stateCopyTail: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte(";\n"),
							},
							DataHandler: dhSecurityNil,
						},
					},
				},
				stateTableValues: {
					NextStates: []fsm.NextState{
						{
							Name: stateInit,
							Switch: fsm.Switch{
								Trigger: []byte("\\."),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{'\n'},
								},
								Escape: false,
							},
							DataHandler: dhSecurityNil,
						},
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte{'\t'},
							},
							DataHandler: dhValue,
						},
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte{'\n'},
							},
							DataHandler: dhValueEnd,
						},
					},
				},
			},
		},
	)
}

func columnType(key string) relfilter.ColumnType {
	t, b := typeKeys[key]
	if b == false {
		return relfilter.ColumnTypeNone
	}
	return t
}
