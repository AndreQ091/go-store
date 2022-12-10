package transact

import (
	"fmt"
	"go-store/core"
	"os"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))
	case "postgres":
		return NewPostgresTransactionLogger(PostgresDBParams{
			dbName:   os.Getenv("TLOG_DB_NAME"),
			host:     os.Getenv("TLOG_DB_HOST"),
			user:     os.Getenv("TLOG_DB_USER"),
			password: os.Getenv("TLOG_DB_PASSWORD"),
		})
	case "":
		return nil, fmt.Errorf("transaction logger type not defined")
	default:
		return nil, fmt.Errorf("no such transaction logger %s", logger)
	}

}
