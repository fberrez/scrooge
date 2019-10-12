package backend

import (
	"database/sql"
	"fmt"

	"github.com/ghodss/yaml"
	// lib/pq is used as a driver.
	_ "github.com/lib/pq"
	"github.com/ovh/configstore"
)

// Backend represents the instance which manages
// the database.
type Backend struct {
	db *sql.DB
}

const (
	databaseKey = "database"
)

// New returns a new backend struct mainly containing a
// database connection, according to the read config file.
func New() (*Backend, error) {
	// gets config
	item, err := configstore.GetItemValue(databaseKey)
	if err != nil {
		return nil, err
	}

	// unmarshals config
	var cfg Config
	if err := yaml.Unmarshal([]byte(item), &cfg); err != nil {
		return nil, err
	}

	// checks config validity
	if err := cfg.Valid(); err != nil {
		return nil, err
	}

	// initializes connection uri and makes connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DBName, cfg.sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// checks if the connection is still alive.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Backend{
		db: db,
	}, nil
}
