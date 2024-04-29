package db

import (
	"encoding/base64"
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

	// TODO: Also set the logger?
	// Logger            Logger
	// LogLevel          LogLevel
	return pgxpool.ParseConfig(connStr)
}
