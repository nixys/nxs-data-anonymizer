package pgsql_anonymize

import fsm "github.com/nixys/nxs-go-fsm"

var (
	stateCreateTableName = fsm.StateName("creat table name")
	stateCreateTableTail = fsm.StateName("creat table tail")
	stateCopySearch      = fsm.StateName("copy search")
	stateTableName       = fsm.StateName("table name")
	stateFieldName       = fsm.StateName("field name")
	stateCopyTail        = fsm.StateName("copy tail")
	stateTableValues     = fsm.StateName("table values")
)
