Cancel requests outage pgbouncer
================================

Usage:
```
$ go get github.com/chobostar/pgbouncer-cancel-request
$ cd ~/go/src/github.com/chobostar/pgbouncer-cancel-request
$ docker-compose up -d
$ make run
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

