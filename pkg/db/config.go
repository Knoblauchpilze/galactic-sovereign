package db

import (
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

const postgresqlConnectionStringTemplate = "postgresql://${user}:${password}@${host}:${port}/${dbname}"

func (c Config) toConnPoolConfig() *pgxpool.Config {
	// https://stackoverflow.com/questions/3582552/what-is-the-format-for-the-postgresql-connection-string-url
	connStr := postgresqlConnectionStringTemplate
	connStr = strings.ReplaceAll(connStr, "${user}", c.User)
	// TODO: URL encode
	connStr = strings.ReplaceAll(connStr, "${password}", c.Password)
	connStr = strings.ReplaceAll(connStr, "${host}", c.Host)
	connStr = strings.ReplaceAll(connStr, "${port}", strconv.Itoa(int(c.Port)))
	connStr = strings.ReplaceAll(connStr, "${dbname}", c.Name)

	// TODO: Handle error
	// TODO: Also set the logger?
	// Logger            Logger
	// LogLevel          LogLevel
	fmt.Printf("str: %s\n", connStr)
	conn, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}

	return conn
}
