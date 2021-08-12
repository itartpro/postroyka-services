module main

go 1.16

require (
	go.mods/dbops v0.0.0
	go.mods/grpcc v0.0.0
	go.mods/hashing v0.0.0
	google.golang.org/grpc v1.36.1
)

replace (
	go.mods/dbops => ./dbops
	go.mods/grpcc => ../shared/grpcc
	go.mods/hashing => ./hashing
)
