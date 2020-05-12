Cancel requests outage pgbouncer
================================

Usage:
```
$ go get github.com/chobostar/pgbouncer-cancel-request
$ cd ~/go/src/github.com/chobostar/pgbouncer-cancel-request
$ docker-compose up -d
$ make run_pgbouncer
```

Check `used_clients`:
```
$ psql -h localhost -p 6432 -U pgbouncer -c "show lists" | grep 'used_clients'
```

Ensure that connections are not available:
```
$ psql -h localhost -p 6432 -U postgres -d db
psql: ERROR:  no more connections allowed (max_client_conn)
```

#### Odyssey

```
$ make run_odyssey
```

and it's okay

## Check killed client

1. Start 2 queries:
- first:
```
$ PGAPPNAME=pgbouncer psql -U postgres -p 6432 -h localhost -d db -c "select pg_sleep(3600)"
```
- second:
```
$ PGAPPNAME=odyssey psql -U postgres -p 6532 -h localhost -d db -c "select pg_sleep(3600)"
```

2. Kill both
```
$ kill $(ps aux | grep '[p]sql' | awk '{print $2}')
```

3. Check active backends:
```
$ psql -U postgres -h localhost -c "select application_name from pg_stat_activity where state != 'idle' and pid != pg_backend_pid()"

 application_name 
------------------
 pgbouncer
(1 row)
```