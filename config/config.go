package config

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	TypeOracleSQLEvolutional = "oracle-sql-evolutional"
	TypeOracleSQLDeferred    = "oracle-sql-deferred"
	TypeOracleSQLEventual    = "oracle-sql-eventual"
	TypeMicrosoftNetWeb      = "microsoft-net-web"
	TypeMicrosoftNetWebCore  = "microsoft-net-web-core"
	TypeMicrosoftNetLib      = "microsoft-net-lib"
	TypeGoWeb                = "go-web"
	TypeGoCommand            = "go-command"
)

// Config es la configuración del proyecto
type Config struct {
	SystemID      string   `yaml:"system-id,omitempty"`
	ApplicationID string   `yaml:"application-id,omitempty"`
	Type          string   `yaml:"type,omitempty"`
	From          []string `yaml:"from,omitempty"`
	TrackerID     string   `yaml:"tracker-id,omitempty"`
	IssueID       string   `yaml:"issue-id,omitempty"`
	Package       Package  `yaml:"package,omitempty"`
	Nexus         Nexus    `yaml:"nexus,omitempty"`
}

type Package struct {
	Format string `yaml:"format,omitempty"`
}

type Nexus struct {
	URL string `yaml:"url,omitempty"`
}

func (c *Config) Validate() error {

	if c.Type == "" {
		return errors.Errorf(`type is required in configuration`)
	}

	if !strings.HasPrefix(c.ApplicationID, c.SystemID+"-") {
		return errors.Errorf(`system-id ("%s-") must be a prefix of application-id (found "%s")`, c.SystemID, c.ApplicationID)
	}

	if c.Type == TypeOracleSQLEventual {
		if c.TrackerID == "" {
			return errors.Errorf(`tracker-id is required in oracle sql eventual configuration`)
		}
	} else {
		if c.IssueID != "" {
			return errors.Errorf(`issue-id must not be defined in configuration`)
		}
		if c.TrackerID != "" {
			return errors.Errorf(`tracker-id must not be defined in configuration`)
		}
	}

	if c.Type != TypeOracleSQLEvolutional {
		if len(c.From) != 0 {
			return errors.Errorf(`form must not be defined in configuration`)
		}
	}

	return nil
}

// Read carga una configuración de proyecto
func Read(r io.Reader) (*Config, error) {
	c := &Config{}
	err := yaml.NewDecoder(r).Decode(c)
	if err != nil {
		return c, errors.Wrapf(err, "decondig project config")
	}
	return c, c.Validate()
}

// ReadFile carga una configuración de proyecto desde un archivo
func ReadFile(location string) (*Config, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, errors.Wrapf(err, "opening '%s'", location)
	}
	defer f.Close()
	c, err := Read(f)
	if err != nil {
		return nil, errors.Wrapf(err, "reading project configuration")
	}
	return c, nil
}
