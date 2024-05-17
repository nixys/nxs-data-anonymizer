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
	Policy     SecurityPolicySettings
	Exceptions SecurityExceptionsSettings
}

type SecurityPolicySettings struct {
	Tables  misc.SecurityPolicyTablesType
	Columns misc.SecurityPolicyColumnsType
}

type SecurityExceptionsSettings struct {
	Tables  map[string]any
	Columns map[string]any
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

	policy     securityPolicyCtx
	exceptions securityExceptionsCtx
}

type securityPolicyCtx struct {
	tables  misc.SecurityPolicyTablesType
	columns misc.SecurityPolicyColumnsType
}

type securityExceptionsCtx struct {
	tables  map[string]any
	columns map[string]any
}

var typeKeys = map[string]relfilter.ColumnType{

	// Special
	"generated": relfilter.ColumnTypeNone,

	// Strings
	"char":       relfilter.ColumnTypeString,
	"varchar":    relfilter.ColumnTypeString,
	"tinytext":   relfilter.ColumnTypeString,
	"text":       relfilter.ColumnTypeString,
	"mediumtext": relfilter.ColumnTypeString,
	"longtext":   relfilter.ColumnTypeString,
	"enum":       relfilter.ColumnTypeString,
	"set":        relfilter.ColumnTypeString,
	"date":       relfilter.ColumnTypeString,
	"datetime":   relfilter.ColumnTypeString,
	"timestamp":  relfilter.ColumnTypeString,
	"time":       relfilter.ColumnTypeString,
	"year":       relfilter.ColumnTypeString,
	"json":       relfilter.ColumnTypeString,

	// Numeric
	"bit":              relfilter.ColumnTypeNum,
	"bool":             relfilter.ColumnTypeNum,
	"boolean":          relfilter.ColumnTypeNum,
	"tinyint":          relfilter.ColumnTypeNum,
	"smallint":         relfilter.ColumnTypeNum,
	"mediumint":        relfilter.ColumnTypeNum,
	"int":              relfilter.ColumnTypeNum,
	"integer":          relfilter.ColumnTypeNum,
	"bigint":           relfilter.ColumnTypeNum,
	"float":            relfilter.ColumnTypeNum,
	"double":           relfilter.ColumnTypeNum,
	"double precision": relfilter.ColumnTypeNum,
	"decimal":          relfilter.ColumnTypeNum,
	"dec":              relfilter.ColumnTypeNum,

	// Binary
	"binary":     relfilter.ColumnTypeBinary,
	"varbinary":  relfilter.ColumnTypeBinary,
	"tinyblob":   relfilter.ColumnTypeBinary,
	"blob":       relfilter.ColumnTypeBinary,
	"mediumblob": relfilter.ColumnTypeBinary,
	"longblob":   relfilter.ColumnTypeBinary,
}

func userCtxInit(s InitSettings) *userCtx {
	return &userCtx{
		filter: relfilter.Init(s.Rules),
		security: securityCtx{
			policy: securityPolicyCtx{
				tables:  s.Security.Policy.Tables,
				columns: s.Security.Policy.Columns,
			},
			exceptions: securityExceptionsCtx{
				tables:  s.Security.Exceptions.Tables,
				columns: s.Security.Exceptions.Columns,
			},
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
							DataHandler: dhSecurityCreateTable,
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
							DataHandler: dhSecurityCreateTableName,
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
							DataHandler: dhSecurityCreateTableName,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
						},
						{
							Name: stateFieldsDescriptionName,
							Switch: fsm.Switch{
								Trigger: []byte("`"),
							},
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
						},
						{
							Name: statefFieldsDescriptionBlockEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
								},
							},
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityCreateTable,
						},
						{
							Name: stateInsertInto,
							Switch: fsm.Switch{
								Trigger: []byte("INSERT"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{' '},
								},
							},
							DataHandler: dhSecurityNil,
						},
					},
				},

				stateInsertInto: {
					NextStates: []fsm.NextState{
						{
							Name: stateInsertIntoTableNameSearch,
							Switch: fsm.Switch{
								Trigger: []byte("INTO"),
								Delimiters: fsm.Delimiters{
									L: []byte{' '},
									R: []byte{' '},
								},
							},
							DataHandler: dhSecurityNil,
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
							DataHandler: dhSecurityNil,
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
