#!/bin/bash

#set -v

SCRIPT_DIR=$(dirname "$(realpath "$0")")
echo "DIR: $SCRIPT_DIR"
echo "PATH: $PATH"
if ! [ -x "$(command -v benchstat)" ]; then
  go install golang.org/x/perf/cmd/benchstat@latest
fi


cd "$SCRIPT_DIR/.." && go test -v -run TestRollbackBenchmark ./instr/

cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5s . > old.out

sampling_rate_list=(0 0.0001 0.001 0.01 0.1 1)

for i in "${!sampling_rate_list[@]}";
do
  echo "### 采样率 ${sampling_rate_list[$i]}"
  out_file=new${i}.out
  cd "$SCRIPT_DIR/.." && go test -v -run TestUpdateBenchmark ./instr/ "${sampling_rate_list[$i]}"
  cd "$SCRIPT_DIR/" && go test -v -run ^$ -bench=. -benchmem -benchtime=5s . > "$out_file"
  cat old.out "$out_file"
  benchstat old.out "$out_file"
done
