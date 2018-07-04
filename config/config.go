package config

import (
	"io"
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
}

func (c *Config) Validate() error {
	if !strings.HasPrefix(c.ApplicationID, c.SystemID+"-") {
		return errors.Errorf(`system-id ("%s-") must be a prefix of application-id (found "%s")`, c.SystemID, c.ApplicationID)
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
