package pgsql_anonymize

import (
	"context"
	"example/modules/filters/relfilter"
	"io"

	fsm "github.com/nixys/nxs-go-fsm"
)

func Init(ctx context.Context, r io.Reader, rules relfilter.Rules) io.Reader {

	return fsm.Init(
		r,
		fsm.Description{
			Ctx:       ctx,
			UserCtx:   relfilter.Init(rules),
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
