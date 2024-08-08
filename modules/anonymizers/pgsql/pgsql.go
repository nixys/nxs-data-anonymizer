package pgsql_anonymize

import (
	"context"
	"fmt"
	"io"

	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"

	fsm "github.com/nixys/nxs-go-fsm"
)

type PgSQL struct {
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
	filter   *relfilter.Filter
	tn       *string
	security securityCtx
	tables   map[string]map[string]string
}

type securityCtx struct {
	tmpBuf []byte
	isSkip bool

	tablesPolicy    misc.SecurityPolicyTablesType
	tableExceptions map[string]any
}

func userCtxInit(s InitOpts) (*userCtx, error) {

	f, err := relfilter.Init(
		relfilter.InitOpts{
			Variables:        s.Variables,
			Link:             s.Link,
			TableRules:       s.Rules.TableRules,
			DefaultRules:     s.Rules.DefaultRules,
			ExceptionColumns: s.Rules.ExceptionColumns,
			TypeRuleCustom:   s.Rules.TypeRuleCustom,
			TypeRuleDefault:  typeRuleDefault,
			ColumnsPolicy:    s.Security.ColumnsPolicy,
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
		tables: make(map[string]map[string]string),
	}, nil
}

func Init(r io.Reader, s InitOpts) (*PgSQL, error) {

	uctx, err := userCtxInit(s)
	if err != nil {
		return nil, fmt.Errorf("pgsql anonymizer init: %w", err)
	}

	return &PgSQL{
		uctx:         uctx,
		sourceReader: r,
	}, nil
}

func (p *PgSQL) Run(ctx context.Context, w io.Writer) error {

	ar := fsm.Init(
		p.sourceReader,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   p.uctx,
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

	_, err := io.Copy(w, ar)
	if err != nil {
		return fmt.Errorf("pgsql anonymizer run: %w", err)
	}

	return nil
}
