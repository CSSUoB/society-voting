package config

import (
	"fmt"
	"git.tdpain.net/pkg/cfger"
	"log/slog"
	"os"
	"sync"
)

type HTTP struct {
	Host string
	Port int
}

func (h *HTTP) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type Guild struct {
	SessionToken string
	SocietyID    string
}

type Database struct {
	DSN string
}

type Platform struct {
	SocietyName string
	AdminToken  string
}

type Config struct {
	Debug    bool
	HTTP     *HTTP
	Guild    *Guild
	Database *Database
	Platform *Platform
}

var (
	conf     *Config
	loadOnce = new(sync.Once)
)

func Get() *Config {
	var outerErr error
	loadOnce.Do(func() {
		cl := cfger.New()
		if err := cl.Load("config.yml"); err != nil {
			outerErr = err
			return
		}

		conf = &Config{
			Debug: cl.WithDefault("debug", false).AsBool(),
			HTTP: &HTTP{
				Host: cl.WithDefault("http.host", "127.0.0.1").AsString(),
				Port: cl.WithDefault("http.port", 8080).AsInt(),
			},
			Guild: &Guild{
				SessionToken: cl.Required("guild.sessionToken").AsString(),
				SocietyID:    cl.Required("guild.societyID").AsString(),
			},
			Database: &Database{
				DSN: cl.WithDefault("database.dsn", "voting.sqlite3.db").AsString(),
			},
			Platform: &Platform{
				SocietyName: cl.WithDefault("platform.societyName", "Society").AsString(),
				AdminToken:  cl.Required("platform.adminToken").AsString(),
			},
		}
	})

	if outerErr != nil {
		slog.Error("fatal error when loading configuration", "err", outerErr)
		os.Exit(1)
	}

	return conf
}
