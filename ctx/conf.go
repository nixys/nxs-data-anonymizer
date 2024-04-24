package ctx

import (
	"fmt"

	"github.com/nixys/nxs-data-anonymizer/misc"
	conf "github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	LogFile  string `conf:"logfile" conf_extraopts:"default=stderr"`
	LogLevel string `conf:"loglevel" conf_extraopts:"default=info"`

	Progress progressConf          `conf:"progress"`
	Filters  map[string]filterConf `conf:"filters"`
	MySQL    *mysqlConf            `conf:"mysql"`
}

type progressConf struct {
	Rhythm   string `conf:"rhythm" conf_extraopts:"default=0s"`
	Humanize bool   `conf:"humanize"`
}

type filterConf struct {
	Columns map[string]columnFilterConf `conf:"columns"`
}

type columnFilterConf struct {
	Type   string `conf:"type" conf_extraopts:"default=template"`
	Value  string `conf:"value" conf_extraopts:"required"`
	Unique bool   `conf:"unique"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"required"`
	Port     int    `conf:"port" conf_extraopts:"required"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	err := conf.Load(&c, conf.Settings{
		ConfPath:    confPath,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, err
	}

	for _, f := range c.Filters {
		for _, cf := range f.Columns {
			if misc.ValueTypeFromString(cf.Type) == misc.ValueTypeUnknown {
				return c, fmt.Errorf("conf read: unknown filter type")
			}
		}
	}

	return c, nil
}
