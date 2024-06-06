package pgsql_anonymize

import fsm "github.com/nixys/nxs-go-fsm"

var (
	stateInit            = fsm.StateName("init")
	stateCreateTableName = fsm.StateName("create table name")
	stateCreateTableTail = fsm.StateName("create table tail")
	stateTableName       = fsm.StateName("table name")
	stateFieldName       = fsm.StateName("field name")
	stateCopyTail        = fsm.StateName("copy tail")
	stateTableValues     = fsm.StateName("table values")
)
