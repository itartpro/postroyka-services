module main

go 1.17

replace go.mods/grpcc => ../shared/grpcc

require (
	github.com/georgysavva/scany v0.2.8
	github.com/jackc/pgx/v4 v4.11.0
	go.mods/grpcc v0.0.0
	google.golang.org/grpc v1.38.0
)
