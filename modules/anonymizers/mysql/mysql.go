package mysql_anonymize

import (
	"context"
	"fmt"
	"io"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"

	fsm "github.com/nixys/nxs-go-fsm"
)

type MySQL struct {
	uctx         *userCtx
	sourceReader io.Reader
}

type InitOpts struct {
	Variables map[string]relfilter.VariableRuleOpts
	Security  SecurityOpts
	Rules     RulesOpts
	Link      []relfilter.LinkOpts
}

type RulesOpts struct {
	TableRules       map[string]map[string]relfilter.ColumnRuleOpts
	DefaultRules     map[string]relfilter.ColumnRuleOpts
	ExceptionColumns []string
	TypeRuleCustom   []relfilter.TypeRuleOpts
}

type SecurityOpts struct {
	TablesPolicy    misc.SecurityPolicyTablesType
	ColumnsPolicy   misc.SecurityPolicyColumnsType
	TableExceptions []string
}

type userCtx struct {
	filter     *relfilter.Filter
	columnName string
	security   securityCtx
	tables     map[string]map[string]columnType
}

type securityCtx struct {
	tmpBuf []byte
	isSkip bool

	tablesPolicy    misc.SecurityPolicyTablesType
	tableExceptions map[string]any
}

type columnType string

const (
	columnTypeNone   columnType = "none"
	columnTypeString columnType = "string"
	columnTypeNum    columnType = "numeric"
	columnTypeBinary columnType = "binary"
)

func (c columnType) String() string {
	return string(c)
}

var typeKeys = map[string]columnType{

	// Strings
	"CHAR":       columnTypeString,
	"VARCHAR":    columnTypeString,
	"TINYTEXT":   columnTypeString,
	"TEXT":       columnTypeString,
	"MEDIUMTEXT": columnTypeString,
	"LONGTEXT":   columnTypeString,
	"ENUM":       columnTypeString,
	"SET":        columnTypeString,
	"DATE":       columnTypeString,
	"DATETIME":   columnTypeString,
	"TIMESTAMP":  columnTypeString,
	"TIME":       columnTypeString,
	"YEAR":       columnTypeString,
	"JSON":       columnTypeString,

	// Numeric
	"BIT":              columnTypeNum,
	"BOOL":             columnTypeNum,
	"BOOLEAN":          columnTypeNum,
	"TINYINT":          columnTypeNum,
	"SMALLINT":         columnTypeNum,
	"MEDIUMINT":        columnTypeNum,
	"INT":              columnTypeNum,
	"INTEGER":          columnTypeNum,
	"BIGINT":           columnTypeNum,
	"FLOAT":            columnTypeNum,
	"DOUBLE":           columnTypeNum,
	"DOUBLE precision": columnTypeNum,
	"DECIMAL":          columnTypeNum,
	"DEC":              columnTypeNum,

	// Binary
	"BINARY":     columnTypeBinary,
	"VARBINARY":  columnTypeBinary,
	"TINYBLOB":   columnTypeBinary,
	"BLOB":       columnTypeBinary,
	"MEDIUMBLOB": columnTypeBinary,
	"LONGBLOB":   columnTypeBinary,
}

func userCtxInit(s InitOpts) (*userCtx, error) {

	trc := []relfilter.TypeRuleOpts{}
	trd := []relfilter.TypeRuleOpts{}
	if s.Security.ColumnsPolicy == misc.SecurityPolicyColumnsRandomize {
		trc = s.Rules.TypeRuleCustom
		trd = typeRuleDefault
	}

	f, err := relfilter.Init(
		relfilter.InitOpts{
			Variables:        s.Variables,
			Link:             s.Link,
			TableRules:       s.Rules.TableRules,
			DefaultRules:     s.Rules.DefaultRules,
			ExceptionColumns: s.Rules.ExceptionColumns,
			TypeRuleCustom:   trc,
			TypeRuleDefault:  trd,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("user ctx init: %w", err)
	}

	return &userCtx{
		filter: f,
		security: securityCtx{
			tablesPolicy: s.Security.TablesPolicy,
			tableExceptions: func() map[string]any {
				excs := make(map[string]any)
				for _, e := range s.Security.TableExceptions {
					excs[e] = nil
				}
				return excs
			}(),
		},
		tables: make(map[string]map[string]columnType),
	}, nil
}

func Init(r io.Reader, s InitOpts) (*MySQL, error) {

	uctx, err := userCtxInit(s)
	if err != nil {
		return nil, fmt.Errorf("mysql anonymizer init: %w", err)
	}

	return &MySQL{
		uctx:         uctx,
		sourceReader: r,
	}, nil
}

func (m *MySQL) Run(ctx context.Context, w io.Writer) error {

	ar := fsm.Init(
		m.sourceReader,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   m.uctx,
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
					NextStates: []fsm.NextState{
						{
							Name: stateFieldsDescriptionBlock,
							Switch: fsm.Switch{
								Trigger: []byte(","),
								Delimiters: fsm.Delimiters{
									R: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableColumnAdd,
						},
						{
							Name: statefFieldsDescriptionBlockEnd,
							Switch: fsm.Switch{
								Trigger: []byte(")"),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
								},
							},
							DataHandler: dhCreateTableColumnAdd,
						},
					},
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

	_, err := io.Copy(w, ar)
	if err != nil {
		return fmt.Errorf("mysql anonymizer run: %w", err)
	}

	return nil
}
