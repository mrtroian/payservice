# !/bin/sh

go test $(pwd)/internal/api/*.go
go test $(pwd)/internal/client/*.go
