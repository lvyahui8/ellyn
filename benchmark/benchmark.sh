#!/bin/sh

set -v

SCRIPT_DIR=$(dirname "$(realpath "$0")")
echo "DIR: $SCRIPT_DIR"
echo "PATH: $PATH"
if ! [ -x "$(command -v benchstat)" ]; then
  go install golang.org/x/perf/cmd/benchstat@latest
fi


cd "$SCRIPT_DIR/.." && go test -v -run TestRollbackBenchmark ./ellyn_ast/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5x . > old.out

cd "$SCRIPT_DIR/.." && go test -v -run TestUpdateBenchmark ./ellyn_ast/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5x . > new.out

benchstat old.out new.out