package ctx

import (
	"example/ds/mysql"
	"example/modules/filters/relfilter"
	"fmt"

	appctx "github.com/nixys/nxs-go-appctx/v2"
)

// Ctx defines application custom context
type Ctx struct {
	Conf  confOpts
	Args  *Args
	Rules relfilter.Rules
	MySQL *mysql.MySQL
}

// Init initiates application custom context
func (c *Ctx) Init(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	a := opts.Args.(*Args)

	// Read config file
	conf, err := confRead(opts.Config)
	if err != nil {
		return appctx.CfgData{}, err
	}

	// Set application context
	c.Conf = conf
	c.Args = a

	// Connect to MySQL if necessary
	if c.Conf.MySQL != nil {
		m, err := mysql.Connect(mysql.Settings{
			Host:     c.Conf.MySQL.Host,
			Port:     c.Conf.MySQL.Port,
			Database: c.Conf.MySQL.DB,
			User:     c.Conf.MySQL.User,
			Password: c.Conf.MySQL.Password,
		})
		if err != nil {
			return appctx.CfgData{}, err
		}
		c.MySQL = &m
	} else {
		if a.Cleanup == true {
			return appctx.CfgData{}, fmt.Errorf("destination database clean up was requested but connection to database doesn't specified")
		}
	}

	c.Rules.Tables = make(map[string]relfilter.TableRules)

	for t, f := range c.Conf.Filters {

		c.Rules.Tables[t] = relfilter.TableRules{
			Columns: func() map[string]relfilter.ColumnRule {
				cc := make(map[string]relfilter.ColumnRule)
				for c, cf := range f.Columns {
					cc[c] = relfilter.ColumnRule{
						Value:  cf.Value,
						Unique: cf.Unique,
					}
				}
				return cc
			}(),
		}
	}

	return appctx.CfgData{
		LogFile:  c.Conf.LogFile,
		LogLevel: c.Conf.LogLevel,
	}, nil
}

// Reload reloads application custom context
func (c *Ctx) Reload(opts appctx.CustomContextFuncOpts) (appctx.CfgData, error) {

	opts.Log.Debug("reloading context")

	return c.Init(opts)
}

// Free frees application custom context
func (c *Ctx) Free(opts appctx.CustomContextFuncOpts) int {

	opts.Log.Debug("freeing context")

	return 0
}
