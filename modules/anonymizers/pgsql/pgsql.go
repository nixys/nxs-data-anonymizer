package pgsql_anonymize

import (
	"context"
	"io"

	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"

	fsm "github.com/nixys/nxs-go-fsm"
)

type InitSettings struct {
	Rules relfilter.Rules
}

type userCtx struct {
	filter *relfilter.Filter

	tn     *string
	tables map[string]map[string]relfilter.ColumnType
}

var typeKeys = map[string]relfilter.ColumnType{

	// Strings
	"character": relfilter.ColumnTypeString,

	// Numeric
	"integer": relfilter.ColumnTypeNum,
}

func userCtxInit(s InitSettings) *userCtx {
	return &userCtx{
		filter: relfilter.Init(s.Rules),
		tables: make(map[string]map[string]relfilter.ColumnType),
	}
}

func Init(ctx context.Context, r io.Reader, s InitSettings) io.Reader {

	return fsm.Init(
		r,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   userCtxInit(s),
			InitState: stateCopySearch,
			States: map[fsm.StateName]fsm.State{

				stateCopySearch: {
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
							DataHandler: nil,
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
							Name: stateCopySearch,
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
							DataHandler: nil,
						},
					},
				},
				stateTableValues: {
					NextStates: []fsm.NextState{
						{
							Name: stateCopySearch,
							Switch: fsm.Switch{
								Trigger: []byte("\\."),
								Delimiters: fsm.Delimiters{
									L: []byte{'\n'},
									R: []byte{'\n'},
								},
								Escape: false,
							},
							DataHandler: nil,
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
