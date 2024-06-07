package mysql_anonymize

import (
	"context"
	"io"
	"strings"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"

	fsm "github.com/nixys/nxs-go-fsm"
)

type InitSettings struct {
	Security SecuritySettings
	Rules    relfilter.Rules
}

type SecuritySettings struct {
	TablePolicy     misc.SecurityPolicyTablesType
	TableExceptions map[string]any
}

type userCtx struct {
	filter *relfilter.Filter
	column userColumnCtx

	security securityCtx
}

type userColumnCtx struct {
	name       string
	columnType relfilter.ColumnType
	isSkip     bool
}

type securityCtx struct {
	tmpBuf []byte
	isSkip bool

	tablePolicy     misc.SecurityPolicyTablesType
	tableExceptions map[string]any
}

const (
	columnTypeString relfilter.ColumnType = "string"
	columnTypeNum    relfilter.ColumnType = "numeric"
	columnTypeBinary relfilter.ColumnType = "binary"
)

var typeKeys = map[string]relfilter.ColumnType{

	// Special
	"generated": relfilter.ColumnTypeNone,

	// Strings
	"char":       columnTypeString,
	"varchar":    columnTypeString,
	"tinytext":   columnTypeString,
	"text":       columnTypeString,
	"mediumtext": columnTypeString,
	"longtext":   columnTypeString,
	"enum":       columnTypeString,
	"set":        columnTypeString,
	"date":       columnTypeString,
	"datetime":   columnTypeString,
	"timestamp":  columnTypeString,
	"time":       columnTypeString,
	"year":       columnTypeString,
	"json":       columnTypeString,

	// Numeric
	"bit":              columnTypeNum,
	"bool":             columnTypeNum,
	"boolean":          columnTypeNum,
	"tinyint":          columnTypeNum,
	"smallint":         columnTypeNum,
	"mediumint":        columnTypeNum,
	"int":              columnTypeNum,
	"integer":          columnTypeNum,
	"bigint":           columnTypeNum,
	"float":            columnTypeNum,
	"double":           columnTypeNum,
	"double precision": columnTypeNum,
	"decimal":          columnTypeNum,
	"dec":              columnTypeNum,

	// Binary
	"binary":     columnTypeBinary,
	"varbinary":  columnTypeBinary,
	"tinyblob":   columnTypeBinary,
	"blob":       columnTypeBinary,
	"mediumblob": columnTypeBinary,
	"longblob":   columnTypeBinary,
}

var RandomizeTypesDefault = map[relfilter.ColumnType]relfilter.ColumnRule{
	columnTypeBinary: {
		Type:   misc.ValueTypeTemplate,
		Value:  "cmFuZG9taXplZCBiaW5hcnkgZGF0YQo=",
		Unique: false,
	},
	columnTypeNum: {
		Type:   misc.ValueTypeTemplate,
		Value:  "0",
		Unique: false,
	},
	columnTypeString: {
		Type:   misc.ValueTypeTemplate,
		Value:  "randomized string data",
		Unique: false,
	},
}

func userCtxInit(s InitSettings) *userCtx {
	return &userCtx{
		filter: relfilter.Init(s.Rules),
		security: securityCtx{
			tablePolicy:     s.Security.TablePolicy,
			tableExceptions: s.Security.TableExceptions,
		},
	}
}

func Init(ctx context.Context, r io.Reader, s InitSettings) io.Reader {

	return fsm.Init(
		r,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   userCtxInit(s),
			InitState: stateCreateSearch,
			States: map[fsm.StateName]fsm.State{

				stateCreateSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateCreateTableSearch,
							Switch: fsm.Switch{
								Trigger: []byte("CREATE"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
					},
				},
				stateCreateTableSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateCreateTableNameSearch,
							Switch: fsm.Switch{
								Trigger: []byte("TABLE"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
					},
				},
				stateCreateTableNameSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateCreateTableName,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: nil,
						},
					},
				},
				stateCreateTableName: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldsDescriptionSearch,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: dhCreateTableName,
						},
					},
				},
				stateFieldsDescriptionSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldsDescriptionBlock,
							Switch: fsm.Switch{
								Trigger: []byte("("),
							},
							DataHandler: nil,
						},
					},
				},
				stateFieldsDescriptionBlock: {
					NextStates: []fsm.NextState{
						{
							// Skip table keys description
							Name: stateFieldDescriptionTailSkip,
							Switch: fsm.Switch{
								Trigger: []byte("KEY"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							// Skip table keys description
							Name: stateFieldDescriptionTailSkip,
							Switch: fsm.Switch{
								Trigger: []byte("PRIMARY"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							// Skip table keys description
							Name: stateFieldDescriptionTailSkip,
							Switch: fsm.Switch{
								Trigger: []byte("UNIQUE"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							// Skip table keys description
							Name: stateFieldDescriptionTailSkip,
							Switch: fsm.Switch{
								Trigger: []byte("CONSTRAINT"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							// Skip table keys description
							Name: stateFieldDescriptionTailSkip,
							Switch: fsm.Switch{
								Trigger: []byte("FOREIGN"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							Name: stateFieldsDescriptionName,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: nil,
						},
					},
				},
				stateFieldDescriptionTailSkip: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldsDescriptionBlock,
							Switch: fsm.Switch{
								Trigger: []byte(","),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: nil,
						},
						{
							Name: statefFieldsDescriptionBlockEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
								},
							},
							DataHandler: nil,
						},
					},
				},
				stateFieldsDescriptionName: {
					NextStates: []fsm.NextState{
						{
							Name: stateFieldsDescriptionNameTail,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: dhCreateTableFieldName,
						},
					},
				},
				stateFieldsDescriptionNameTail: {
					NextStates: func() []fsm.NextState {

						var nss []fsm.NextState

						for t := range typeKeys {
							for i := 0; i < 2; i++ {

								s := t
								if i == 1 {
									s = strings.ToUpper(t)
								}

								nss = append(nss, fsm.NextState{
									Name: stateFieldsDescriptionNameTail,
									Switch: fsm.Switch{
										Trigger: []byte(s),
										Delimiters: fsm.Delimiters{
											L: []byte{' '},
											R: []byte{' ', '(', ',', '\n'},
										},
									},
									DataHandler: dhCreateTableColumnTypeAdd,
								})
							}
						}

						// Additional states
						nss = append(nss, fsm.NextState{
							Name: stateFieldsDescriptionBlock,
							Switch: fsm.Switch{
								Trigger: []byte(","),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableColumnAdd,
						})
						nss = append(nss, fsm.NextState{
							Name: statefFieldsDescriptionBlockEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableColumnAdd,
						})

						return nss
					}(),
				},
				statefFieldsDescriptionBlockEnd: {
					NextStates: []fsm.NextState{
						{
							Name: stateSomeIntermediateState,
							Switch: fsm.Switch{
								Trigger: []byte(";"),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: nil,
						},
					},
				},
				stateSomeIntermediateState: {
					NextStates: []fsm.NextState{
						{
							Name: stateCreateTableSearch,
							Switch: fsm.Switch{
								Trigger: []byte("CREATE"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: nil,
						},
						{
							Name: stateInsertIntoTableNameSearch,
							Switch: fsm.Switch{
								Trigger: []byte("INSERT INTO"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: dhSecurityInsertInto,
						},
					},
				},

				stateInsertIntoTableNameSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateInsertIntoTableName,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: dhSecurityInsertIntoTableNameSearch,
						},
					},
				},
				stateInsertIntoTableName: {
					NextStates: []fsm.NextState{
						{
							Name: stateValuesSearchKeyword,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: dhInsertIntoTableName,
						},
					},
				},
				stateValuesSearchKeyword: {
					NextStates: []fsm.NextState{
						{
							Name: stateValuesSearch,
							Switch: fsm.Switch{
								Trigger: []byte("VALUES"),
								Delimiters: fsm.Delimiters{
									L: []byte{' ', '\n'},
									R: []byte{' ', '\n'},
								},
							},
							DataHandler: dhSecurityNil,
						},
					},
				},
				stateValuesSearch: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte("("),
							},
							DataHandler: fsm.DataHandlerGenericSkipToken,
						},
					},
				},
				stateTableValues: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableValuesString,
							Switch: fsm.Switch{
								Trigger: []byte("'"),
							},
							DataHandler: fsm.DataHandlerGenericVoid,
						},
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte(","),
							},
							DataHandler: dhCreateTableValues,
						},
						{
							Name: stateTableValuesEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
							},
							DataHandler: dhCreateTableValuesEnd,
						},
					},
				},
				stateTableValuesString: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableValuesStringEnd,
							Switch: fsm.Switch{
								Trigger: []byte("'"),
								Escape:  true,
							},
							DataHandler: dhCreateTableValuesString,
						},
					},
				},
				stateTableValuesStringEnd: {
					NextStates: []fsm.NextState{
						{
							Name: stateTableValues,
							Switch: fsm.Switch{
								Trigger: []byte(","),
							},
							DataHandler: fsm.DataHandlerGenericVoid,
						},
						{
							Name: stateTableValuesEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
							},
							DataHandler: dhCreateTableValuesStringEnd,
						},
					},
				},
				stateTableValuesEnd: {
					NextStates: []fsm.NextState{
						{
							Name: stateValuesSearch,
							Switch: fsm.Switch{
								Trigger: []byte(","),
							},
							DataHandler: dhSecurityNil,
						},
						{
							Name: stateSomeIntermediateState,
							Switch: fsm.Switch{
								Trigger: []byte(";"),
							},
							DataHandler: dhSecurityNil,
						},
					},
				},
			},
		},
	)
}
