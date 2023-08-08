package mysql_anonymize

import fsm "github.com/nixys/nxs-go-fsm"

var (
	stateCreateSearch               = fsm.StateName("create search")
	stateCreateTableSearch          = fsm.StateName("create table search")
	stateCreateTableNameSearch      = fsm.StateName("create table name search")
	stateCreateTableName            = fsm.StateName("create table name")
	stateFieldsDescriptionSearch    = fsm.StateName("fields description search")
	stateFieldsDescriptionBlock     = fsm.StateName("fields description block")
	stateFieldsDescriptionName      = fsm.StateName("fields description name")
	stateFieldsDescriptionNameTail  = fsm.StateName("fields description name tail")
	stateFieldDescriptionTailSkip   = fsm.StateName("fields description tail skip")
	statefFieldsDescriptionBlockEnd = fsm.StateName("fields description block end")
	stateSomeIntermediateState      = fsm.StateName("some intermediate state")
	stateInsertInto                 = fsm.StateName("insert into")
	stateInsertIntoTableNameSearch  = fsm.StateName("insert into table name search")
	stateInsertIntoTableName        = fsm.StateName("insert into table name")
	stateValuesSearch               = fsm.StateName("values search")
	stateTableValues                = fsm.StateName("table values")
	stateTableValuesString          = fsm.StateName("table values string")
	stateTableValuesBinary          = fsm.StateName("table values binary")
	stateTableValuesEnd             = fsm.StateName("table values end")
	stateTableValuesStringEnd       = fsm.StateName("table values string end")
	stateValuesSearchKeyword        = fsm.StateName("values search key VALUES")
	stateFieldsGenerated            = fsm.StateName("fields generated search")
)
