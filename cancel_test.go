package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const dockerHost = "localhost"

var (
	dsnPgbouncer190Double = fmt.Sprintf("host=%s port=6435 database=db1 user=postgres sslmode=disable", dockerHost)

	dsnPgbouncer190 = fmt.Sprintf("host=%s port=6434 database=db1 user=postgres sslmode=disable", dockerHost)
	dsnPgbouncer181 = fmt.Sprintf("host=%s port=6433 database=db1 user=postgres sslmode=disable", dockerHost)
	dsnPgbouncer172 = fmt.Sprintf("host=%s port=6432 database=db1 user=postgres sslmode=disable", dockerHost)
)

func TestCancelRequest(t *testing.T) {
	t.Run("double pgbouncer 1.9.0", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer190Double)
		require.Nil(t, err)

		testCancelRequest(t, db)
	})

	t.Run("pgbouncer 1.9.0", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer190)
		require.Nil(t, err)

		testCancelRequest(t, db)
	})

	t.Run("pgbouncer 1.8.1", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer181)
		require.Nil(t, err)

		testCancelRequest(t, db)
	})

	// protocol out of sync after cancel
	t.Run("pgbouncer 1.7.2", func(t *testing.T) {
		db, err := sql.Open("postgres", dsnPgbouncer172)
		require.Nil(t, err)

		require.Panics(t, func() { testCancelRequest(t, db) })
	})
}

func testCancelRequest(t *testing.T, db *sql.DB) {
	const maxConn = 32

	var wg sync.WaitGroup
	var panicsCnt int32

	err := db.Ping()
	require.Nil(t, err)

	wg.Add(maxConn)
	for i := 0; i < maxConn; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					atomic.AddInt32(&panicsCnt, 1)
				}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			rows, err := db.QueryContext(ctx, "select pg_sleep(0.5)")
			if err != nil {
				return
			}

			rows.Close()
		}()
	}

	wg.Wait()

	if panicsCnt > 0 {
		panic(fmt.Sprintf("panics occured: %d", panicsCnt))
	}

	var nonIdleQueriesCnt int
	err = db.QueryRow("SELECT count(1) FROM pg_stat_activity WHERE state <> 'idle'").Scan(&nonIdleQueriesCnt)

	require.Nil(t, err)
	require.Equal(t, 1, nonIdleQueriesCnt)
}
