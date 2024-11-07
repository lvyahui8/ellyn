#!/bin/sh

set -v

SCRIPT_DIR=$(dirname "$(realpath "$0")")
echo "DIR: $SCRIPT_DIR"
echo "PATH: $PATH"
if ! [ -x "$(command -v benchstat)" ]; then
  go install golang.org/x/perf/cmd/benchstat@latest
fi


cd "$SCRIPT_DIR/.." && go test -v -run TestRollbackBenchmark ./instr/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5s . > old.out

cd "$SCRIPT_DIR/.." && go test -v -run TestUpdateBenchmark ./instr/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5s . > new.out

cat old.out new.out
benchstat old.out new.out