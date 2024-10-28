#!/bin/sh

set -v

SCRIPT_DIR=$(dirname "$(realpath "$0")")
echo "DIR: $SCRIPT_DIR"
echo "PATH: $PATH"

go install golang.org/x/perf/cmd/benchstat@latest

cd "$SCRIPT_DIR/.." && go test -v -run TestRollbackBenchmark ./ellyn_ast/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem . > old.out

cd "$SCRIPT_DIR/.." && go test -v -run TestUpdateBenchmark ./ellyn_ast/

cd "$SCRIPT_DIR/" && go test -v -run ^$  -bench=. -benchmem . > new.out

benchstat old.out new.out