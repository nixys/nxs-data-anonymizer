package ctx

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nixys/nxs-data-anonymizer/ds/mysql"
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
	"github.com/sirupsen/logrus"

	appctx "github.com/nixys/nxs-go-appctx/v3"
)

// Ctx defines application custom context
type Ctx struct {
	Log      *logrus.Logger
	Input    io.Reader
	Output   io.Writer
	Rules    relfilter.Rules
	Progress progressCtx
	DB       DBCtx
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

// Init initiates application custom context
func AppCtxInit() (any, error) {

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
		c.Input = os.Stdin
	} else {
		c.Input, err = os.Open(*args.Input)
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

	// Connect to MySQL if necessary
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

	c.Rules.Tables = make(map[string]relfilter.TableRules)

	for t, f := range conf.Filters {

		c.Rules.Tables[t] = relfilter.TableRules{
			Columns: func() map[string]relfilter.ColumnRule {
				cc := make(map[string]relfilter.ColumnRule)
				for c, cf := range f.Columns {
					cc[c] = relfilter.ColumnRule{
						Type:   misc.ValueTypeFromString(cf.Type),
						Value:  cf.Value,
						Unique: cf.Unique,
					}
				}
				return cc
			}(),
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
