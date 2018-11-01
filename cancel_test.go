package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const dockerHost = "localhost"

var (
	dsnPgbouncer190 = fmt.Sprintf("host=%s port=6434 database=db1 user=postgres sslmode=disable", dockerHost)
	dsnPgbouncer181 = fmt.Sprintf("host=%s port=6433 database=db1 user=postgres sslmode=disable", dockerHost)
	dsnPgbouncer172 = fmt.Sprintf("host=%s port=6432 database=db1 user=postgres sslmode=disable", dockerHost)
)

func TestCancelRequest(t *testing.T) {
	t.Run("pgbouncer 1.9.0", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer190)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		testCancelRequest(t, db)
	})

	t.Run("pgbouncer 1.8.1", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer181)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		testCancelRequest(t, db)
	})

	// protocol out of sync after cancel
	t.Run("pgbouncer 1.7.2", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer172)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		testCancelRequest(t, db)
	})
}

func testCancelRequest(t *testing.T, db *sql.DB) {
	const maxConn = 32

	var wg sync.WaitGroup

	db.SetMaxOpenConns(maxConn)

	if err := db.Ping(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < maxConn; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

			rows, err := db.QueryContext(ctx, "select pg_sleep(0.5)")
			if err != nil {
				return
			}

			rows.Close()
		}()
	}

	wg.Wait()
}
