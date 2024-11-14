#!/bin/bash

#set -v

SCRIPT_DIR=$(dirname "$(realpath "$0")")
echo "DIR: $SCRIPT_DIR"
echo "PATH: $PATH"
if ! [ -x "$(command -v benchstat)" ]; then
  go install golang.org/x/perf/cmd/benchstat@latest
fi


cd "$SCRIPT_DIR/.." && go test  -run TestRollbackBenchmark ./instr/

cd "$SCRIPT_DIR/" && go test -run ^$ -bench=. -benchmem -benchtime=5s . > old.out

echo "### 基准"

echo '```text'
cat old.out
echo '```'


sampling_rate_list=(0 0.0001 0.001 0.01 0.1 1)

for i in "${!sampling_rate_list[@]}"; do
  echo "### 采样率 ${sampling_rate_list[$i]}"
  out_file=new${i}.out
  cd "$SCRIPT_DIR/../instr" && go  test -test.paniconexit0 -test.run '^\QTestUpdateBenchmark\E$' ${sampling_rate_list[$i]}
  cd "$SCRIPT_DIR/" && go test  -run ^$ -bench=. -benchmem -benchtime=5s . > "$out_file"
  echo '```text'
  cat "$out_file"
  echo '```'

  echo '性能差异'
  echo '```text'
  benchstat old.out "$out_file"
  echo '```'
done
