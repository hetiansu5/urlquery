#!/bin/bash

curDir=$(cd "$(dirname "$0")"; pwd)

echo "map:"
go run ${curDir}/map.go

echo "\n simple:"
go run ${curDir}/simple.go

echo "\n anonymous:"
go run ${curDir}/struct_anonymous.go

echo "\n converter:"
go run ${curDir}/converter.go

echo "\n option:"
go run ${curDir}/withoption.go

echo "\ngo test:"
go test ${curDir}/../*.go

# benchmark
# go test -bench=. -benchtime=3s -run=none

