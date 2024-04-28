package db

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host                string
	Port                uint16
	Name                string
	User                string
	Password            string
	ConnectionsPoolSize uint
}

// https://github.com/jackc/pgx/blob/60a01d044a5b3f65b9eea866954fdeea1e7d3f00/pgxpool/pool.go#L286
const postgresqlConnectionStringTemplate = "postgresql://${user}:${password}@${host}:${port}/${dbname}?pool_min_conns=${min_connections}"

func (c Config) toConnPoolConfig() (*pgxpool.Config, error) {
	// https://stackoverflow.com/questions/3582552/what-is-the-format-for-the-postgresql-connection-string-url
	connStr := postgresqlConnectionStringTemplate
	connStr = strings.ReplaceAll(connStr, "${user}", c.User)
	// https://stackoverflow.com/questions/58419348/is-there-a-urlencode-function-in-golang
	passwordEncoded := base64.URLEncoding.EncodeToString([]byte(c.Password))
	connStr = strings.ReplaceAll(connStr, "${password}", passwordEncoded)
	connStr = strings.ReplaceAll(connStr, "${host}", c.Host)
	connStr = strings.ReplaceAll(connStr, "${port}", strconv.Itoa(int(c.Port)))
	connStr = strings.ReplaceAll(connStr, "${dbname}", c.Name)
	connStr = strings.ReplaceAll(connStr, "${min_connections}", strconv.Itoa(int(c.ConnectionsPoolSize)))

	var conf *pgxpool.Config
	var parseErr, recoverErr error
	func() {
		defer func() {
			if maybeErr := recover(); maybeErr != nil {
				if err, ok := maybeErr.(error); ok {
					recoverErr = err
				} else {
					recoverErr = fmt.Errorf("%v", maybeErr)
				}
			}
		}()

		conf, parseErr = pgxpool.ParseConfig(connStr)
	}()

	// TODO: Also set the logger?
	// Logger            Logger
	// LogLevel          LogLevel
	fmt.Printf("str: %s\n", connStr)
	err := parseErr
	if parseErr == nil && recoverErr != nil {
		err = recoverErr
	}

	return conf, err
}
