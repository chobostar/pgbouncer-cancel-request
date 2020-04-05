package main

import (
	"context"
	"flag"
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

func test_db(db *sql.DB, q string) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
	result, err := db.ExecContext(ctx, q)
	if err != nil {
		fmt.Println(err.Error())
		return
		//panic(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		return
		//panic(err)
	}
	fmt.Println("Rows affected", rows)
}

func main() {
	var dsn = flag.String("dsn", "postgres://postgres@localhost:6432?dbname=db&sslmode=disable", "PostgreSQL DSN postgres://postgres:password@localhost:5432?sslmode=disable")
	var q = flag.String("q", "SELECT TRUE;", "Query to execute example: \"SELECT TRUE;\"")
	flag.Parse()
	fmt.Println("Using PostgreSQL DSN:", *dsn)
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	var wg sync.WaitGroup
	for e := 0; e < 40; e++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10000; i++ {
				test_db(db, *q)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
