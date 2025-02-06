package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/4nth0/golem/server"
	_ "github.com/lib/pq"
)

type Client struct {
	dbInfos *ConnectionInfos
	db      *sql.DB
}

func NewClient(connection string) *Client {
	parsedConnectionString, err := ParseConnectionString(connection)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", parsedConnectionString.String())
	if err != nil {
		panic(err)
	}
	return &Client{
		db:      db,
		dbInfos: parsedConnectionString,
	}
}

func (f Client) WriteLine(line server.InboundRequest) error {
	sqlStatement := `
INSERT INTO ` + f.dbInfos.DBName + ` (line, created_at)
VALUES ($1, $2)`
	_, err := f.db.Exec(sqlStatement, line, time.Now())
	if err != nil {
		panic(err)
	}
	return nil
}

func (f Client) Close() {
	f.db.Close()
}

type ConnectionInfos struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func ParseConnectionString(connectionString string) (*ConnectionInfos, error) {
	output := &ConnectionInfos{}

	splited := strings.Split(connectionString, " ")
	for _, line := range splited {
		splitedLine := strings.Split(line, "=")

		switch strings.ToLower(splitedLine[0]) {
		case "host":
			output.Host = splitedLine[1]
		case "port":
			port, err := strconv.Atoi(splitedLine[1])
			if err != nil {
				if err != nil {
					panic(err)
				}
			} else {
				output.Port = port
			}
		case "user":
			output.User = splitedLine[1]
		case "password":
			output.Password = splitedLine[1]
		case "dbname":
			output.DBName = splitedLine[1]
		}
	}

	if output.User == "" {
		return nil, errors.New("USER_NOT_SPECIFIED")
	}
	if output.DBName == "" {
		return nil, errors.New("DBNAME_NOT_SPECIFIED")
	}
	if output.Host == "" {
		output.Host = "localhost"
	}
	if output.Port == 0 {
		output.Port = 5432
	}

	return output, nil
}

func (c ConnectionInfos) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
