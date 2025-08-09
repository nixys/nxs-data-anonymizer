// Exemple de modification pour ctx/conf.go
package ctx

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"github.com/nixys/nxs-data-anonymizer/misc"
	conf "github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	LogFile  string `conf:"logfile" conf_extraopts:"default=stderr"`
	LogLevel string `conf:"loglevel" conf_extraopts:"default=info"`

	Progress  progressConf                  `conf:"progress"`
	Filters   map[string]filterConf         `conf:"filters"`
	Link      []linkConf                    `conf:"link"`
	Security  securityConf                  `conf:"security"`
	Variables map[string]variableFilterConf `conf:"variables"`

	// Configuration faker NOUVELLES OPTIONS
	Faker     bool          `conf:"faker"`      // Simple boolean pour activer faker
	FakerConf *fakerConf    `conf:"faker_conf"` // Configuration avancée (optionnelle)
	FakerFile string        `conf:"faker-file"` // Fichier de profils faker
	
	MySQL *mysqlConf `conf:"mysql"`
}

// NOUVELLE structure de configuration faker
type fakerConf struct {
	Enabled bool   `conf:"enabled" conf_extraopts:"default=false"`
	Locale  string `conf:"locale" conf_extraopts:"default=fr_FR"`
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

// NOUVELLE structure pour les profils faker chargés depuis fichier externe
type FakerProfile struct {
	Locale      string                       `yaml:"locale"`
	Description string                       `yaml:"description"`
	Profiles    map[string]interface{}       `yaml:"profiles"`
}

type linkConf struct {
	Rule columnFilterConf    `conf:"rule"`
	With map[string][]string `conf:"with" conf_extraopts:"required"`
}

type variableFilterConf struct {
	Type  string `conf:"type" conf_extraopts:"default=template"`
	Value string `conf:"value" conf_extraopts:"required"`
}

type securityConf struct {
	Policy     securityPolicyConf     `conf:"policy"`
	Exceptions securityExceptionsConf `conf:"exceptions"`
	Defaults   securityDefaultsConf   `conf:"defaults"`
}

type securityPolicyConf struct {
	Tables  string `conf:"tables" conf_extraopts:"default=pass"`
	Columns string `conf:"columns" conf_extraopts:"default=pass"`
}

type securityExceptionsConf struct {
	Tables  []string `conf:"tables"`
	Columns []string `conf:"columns"`
}

type securityDefaultsConf struct {
	Columns map[string]columnFilterConf `conf:"columns"`
	Types   []securityDefaultsTypeConf  `conf:"types"`
}

type securityDefaultsTypeConf struct {
	Regex string           `conf:"regex" conf_extraopts:"required"`
	Rule  columnFilterConf `conf:"rule" conf_extraopts:"required"`
}

type mysqlConf struct {
	Host     string `conf:"host" conf_extraopts:"required"`
	Port     int    `conf:"port" conf_extraopts:"required"`
	DB       string `conf:"db" conf_extraopts:"required"`
	User     string `conf:"user" conf_extraopts:"required"`
	Password string `conf:"password" conf_extraopts:"required"`
}

// NOUVELLE fonction pour charger le profil faker
func loadFakerProfile(filePath string) (*FakerProfile, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read faker profile file: %w", err)
	}

	var profile FakerProfile
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse faker profile file: %w", err)
	}

	return &profile, nil
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

	// NOUVELLE validation pour les types faker
	for _, f := range c.Filters {
		for _, cf := range f.Columns {
			if misc.ValueTypeFromString(cf.Type) == misc.ValueTypeUnknown {
				return c, fmt.Errorf("conf read: unknown column filter type")
			}
		}
	}

	for _, f := range c.Variables {
		if misc.ValueTypeFromString(f.Type) == misc.ValueTypeUnknown {
			return c, fmt.Errorf("conf read: unknown variable filter type")
		}
	}

	if misc.SecurityPolicyTablesTypeFromString(c.Security.Policy.Tables) == misc.SecurityPolicyTablesUnknown {
		return c, fmt.Errorf("conf read: unknown security policy tables type")
	}

	if misc.SecurityPolicyColumnsTypeFromString(c.Security.Policy.Columns) == misc.SecurityPolicyColumnsUnknown {
		return c, fmt.Errorf("conf read: unknown security policy columns type")
	}

	for _, cf := range c.Security.Defaults.Columns {
		if misc.ValueTypeFromString(cf.Type) == misc.ValueTypeUnknown {
			return c, fmt.Errorf("conf read: unknown default filter type")
		}
	}

	// NOUVELLE validation et chargement du profil faker si spécifié
	if c.FakerFile != "" {
		_, err := loadFakerProfile(c.FakerFile)
		if err != nil {
			return c, fmt.Errorf("conf read: failed to load faker profile: %w", err)
		}
	}

	return c, nil
}