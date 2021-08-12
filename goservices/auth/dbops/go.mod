module go.mods/dbops

go 1.16

require (
	github.com/georgysavva/scany v0.2.8
	github.com/jackc/pgx/v4 v4.11.0
	go.mods/hashing v0.0.0
)

replace go.mods/hashing => ../hashing
