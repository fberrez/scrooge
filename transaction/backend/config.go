package backend

import "errors"

// Config is the backend configuration.
type Config struct {
	// DBName is the database name
	DBName string `json:"db_name"`
	// User is the database user
	User string `json:"user"`
	// Password is the database password for the used user
	Password string `json:"password"`
	// Addr is the database address
	Addr string `json:"addr"`
	// Port is the database port
	Port int `json:"port"`
	// SkipVerify is a bool which determines if the backend
	// must skips the ssl mode or not.
	SkipVerify bool `json:"skip_verify,omitempty"`
	sslMode    string
}

// Valid returns an error if the config is not valid.
func (c *Config) Valid() error {
	if c.DBName == "" {
		return errors.New("db name is missing or empty")
	}

	if c.Password == "" {
		return errors.New("password is missing or empty")
	}

	if c.Addr == "" {
		return errors.New("addr is mising or empty")
	}

	if c.Port <= 0 {
		return errors.New("port is missing or empty")
	}

	if c.SkipVerify {
		c.sslMode = "disable"
	} else {
		c.sslMode = "verify-full"
	}

	return nil
}
