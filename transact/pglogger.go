package transact

import (
	"database/sql"
	"fmt"
	"go-store/core"
)

type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- core.Event
	errors <-chan error
	db     *sql.DB
}


func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", config.host, config.dbName, config.user, config.password)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	return logger, nil
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{Type: core.EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{Type: core.EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions
				  (type, key, value)
				  VALUES ($1, $2, $3)`
		for e := range events {
			_, err := l.db.Exec(query, e.Type, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT * from transactions ORDER BY sequense`
		rows, err := l.db.Query(query)

		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		e := core.Event{}

		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.Type, &e.Key, &e.Value)

			if err != nil {
				outError <- fmt.Errorf("reading row error: %w", err)
				return
			}

			outEvent <- e
		}

		err = rows.Err()

		if err != nil {
			outError <- fmt.Errorf("ransaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}