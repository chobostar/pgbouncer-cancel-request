#!/usr/bin/env bash

get_stats() {
  prefix=${1:-pgbouncer}
  port=${2:-6432}

  psql -U pgbouncer -h localhost -p ${port} -d pgbouncer -c "show databases" -xt | cat | awk {'print $1'} > "${prefix}_show_databases.out"
  psql -U pgbouncer -h localhost -p ${port} -d pgbouncer -c "show lists" -t | cat | awk {'print $1'} > "${prefix}_show_lists.out"
  psql -U pgbouncer -h localhost -p ${port} -d pgbouncer -c "show pools" -xt | cat | awk {'print $1'} > "${prefix}_show_pools.out"
  psql -U pgbouncer -h localhost -p ${port} -d pgbouncer -c "show stats" -xt | cat | awk {'print $1'} > "${prefix}_show_stats.out"
}

get_stats pgbouncer 6432
get_stats odyssey 6532

for stat in databases lists pools stats
do
  echo "diff of show ${stat}:"
  diff $(ls *${stat}*.out)
  echo "---"
done
