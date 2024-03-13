package database

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect() (*sqlx.DB, error) {
	c, err := pgx.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	c.Host = os.Getenv("POSTGRES_HOST")

	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	c.Port = uint16(port)

	c.User = os.Getenv("POSTGRES_USER")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Database = os.Getenv("POSTGRES_DATABASE")

	c.TLSConfig = nil
	c.PreferSimpleProtocol = true
	c.RuntimeParams["standard_conforming_strings"] = "on"
	c.DialFunc = (&net.Dialer{
		KeepAlive: 5 * time.Minute,
		Timeout:   1 * time.Second,
	}).DialContext

	return sqlx.NewDb(stdlib.OpenDB(*c), "pgx"), nil
}
