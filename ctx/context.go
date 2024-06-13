package ctx

import (
	"fmt"
	"io"
	"os"
	"time"

	mysql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/mysql"
	pgsql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/pgsql"
	progressreader "github.com/nixys/nxs-data-anonymizer/modules/progress_reader"

	"github.com/nixys/nxs-data-anonymizer/ds/mysql"
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
	"github.com/sirupsen/logrus"

	appctx "github.com/nixys/nxs-go-appctx/v3"
)

// Ctx defines application custom context
type Ctx struct {
	Log        *logrus.Logger
	Output     io.Writer
	Progress   progressCtx
	DB         DBCtx
	Anonymizer misc.Anonymizer
	PR         *progressreader.ProgressReader
}

type DBCtx struct {
	Cleanup bool
	Type    DBType
	MySQL   *mysql.MySQL
}

type DBType string

const (
	DBTypeMySQL DBType = "mysql"
	DBTypePgSQL DBType = "pgsql"
)

type LogFormat string

const (
	LogFormatJSON  LogFormat = "json"
	LogFormatPlain LogFormat = "plain"
)

type progressCtx struct {
	Rhythm   time.Duration
	Humanize bool
}

type SecurityCtx struct {
	TablePolicy     misc.SecurityPolicyTablesType
	TableExceptions map[string]any
}

// Init initiates application custom context
func AppCtxInit() (any, error) {

	var ir io.Reader

	c := &Ctx{}

	args, err := ArgsRead()
	if err != nil {
		return nil, err
	}

	conf, err := confRead(args.ConfigPath)
	if err != nil {
		tmpLogError("ctx init", err)
		return nil, err
	}

	c.Log, err = logInit(conf.LogFile, conf.LogLevel, args.LogFormat)
	if err != nil {
		tmpLogError("ctx init", err)
		return nil, err
	}

	if args.Input == nil {
		ir = os.Stdin
	} else {
		ir, err = os.Open(*args.Input)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init: open input file")
			return nil, err
		}
	}

	if args.Output == nil {
		c.Output = os.Stdout
	} else {
		c.Output, err = os.Create(*args.Output)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init: open output file")
			return nil, err
		}
	}

	c.DB = DBCtx{
		Cleanup: args.Cleanup,
		Type:    args.DBType,
	}

	// DEPRECATED: Connect to MySQL if necessary
	if conf.MySQL != nil {
		m, err := mysql.Connect(mysql.Settings{
			Host:     conf.MySQL.Host,
			Port:     conf.MySQL.Port,
			Database: conf.MySQL.DB,
			User:     conf.MySQL.User,
			Password: conf.MySQL.Password,
		})
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init")
			return nil, err
		}
		c.DB.MySQL = &m
	} else {
		if args.Cleanup == true {
			c.Log.WithFields(logrus.Fields{
				"details": "destination database clean up was requested but connection to database doesn't specified",
			}).Errorf("ctx init")
			return nil, misc.ErrConig
		}
	}

	c.PR = progressreader.Init(ir)

	TableRules := func() map[string]map[string]relfilter.ColumnRuleOpts {
		tables := make(map[string]map[string]relfilter.ColumnRuleOpts)
		for t, cs := range conf.Filters {
			columns := make(map[string]relfilter.ColumnRuleOpts)
			for c, f := range cs.Columns {
				columns[c] = relfilter.ColumnRuleOpts{
					Type:   misc.ValueType(f.Type),
					Value:  f.Value,
					Unique: f.Unique,
				}
			}
			tables[t] = columns
		}
		return tables
	}()

	DefaultRules := func() map[string]relfilter.ColumnRuleOpts {
		cc := make(map[string]relfilter.ColumnRuleOpts)
		for c, cf := range conf.Security.Defaults.Columns {
			cc[c] = relfilter.ColumnRuleOpts{
				Type:   misc.ValueTypeFromString(cf.Type),
				Value:  cf.Value,
				Unique: cf.Unique,
			}
		}
		return cc
	}()

	switch args.DBType {
	case DBTypeMySQL:
		c.Anonymizer, err = mysql_anonymize.Init(
			c.PR,
			mysql_anonymize.InitOpts{
				Security: mysql_anonymize.SecurityOpts{
					TablesPolicy:    misc.SecurityPolicyTablesType(conf.Security.Policy.Tables),
					ColumnsPolicy:   misc.SecurityPolicyColumnsTypeFromString(conf.Security.Policy.Columns),
					TableExceptions: conf.Security.Exceptions.Tables,
				},
				Rules: mysql_anonymize.RulesOpts{
					TableRules:       TableRules,
					DefaultRules:     DefaultRules,
					ExceptionColumns: conf.Security.Exceptions.Columns,
					//TypeRuleCustom: ,
				},
			},
		)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init")
			return nil, err
		}
	case DBTypePgSQL:
		c.Anonymizer, err = pgsql_anonymize.Init(
			c.PR,
			pgsql_anonymize.InitOpts{
				Security: pgsql_anonymize.SecurityOpts{
					TablesPolicy:    misc.SecurityPolicyTablesType(conf.Security.Policy.Tables),
					ColumnsPolicy:   misc.SecurityPolicyColumnsTypeFromString(conf.Security.Policy.Columns),
					TableExceptions: conf.Security.Exceptions.Tables,
				},
				Rules: pgsql_anonymize.RulesOpts{
					TableRules:       TableRules,
					DefaultRules:     DefaultRules,
					ExceptionColumns: conf.Security.Exceptions.Columns,
					//TypeRuleCustom: ,
				},
			},
		)
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("ctx init")
			return nil, err
		}
	}

	// Progress settings
	c.Progress.Humanize = conf.Progress.Humanize

	c.Progress.Rhythm, err = time.ParseDuration(conf.Progress.Rhythm)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"details": err,
		}).Errorf("ctx init")
		return nil, err
	}

	return c, nil
}

func tmpLogError(msg string, err error) {
	l, _ := appctx.DefaultLogInit(os.Stderr, logrus.InfoLevel, &logrus.JSONFormatter{})
	l.WithFields(logrus.Fields{
		"details": err,
	}).Errorf(msg)
}

func logInit(file, level string, ft LogFormat) (*logrus.Logger, error) {

	var (
		f   *os.File
		err error
	)

	switch file {
	case "stdout":
		f = os.Stdout
	case "stderr":
		f = os.Stderr
	default:
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("log init: %w", err)
		}
	}

	// Validate log level
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("log init: %w", err)
	}

	if ft == LogFormatPlain {
		return appctx.DefaultLogInit(f, l, nil)
	}

	return appctx.DefaultLogInit(f, l, &logrus.JSONFormatter{})
}
