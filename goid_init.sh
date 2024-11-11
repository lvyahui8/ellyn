#!/bin/bash

set -e

GOROOT=`go env GOROOT`

RUNTIME_FILE=$GOROOT/src/runtime/runtime2.go
TARGET_FILE=$GOROOT/src/runtime/ellyn_goid.go

updated=$(grep -c -e  'routineCtx\s*uintptr' ${RUNTIME_FILE} || true)

if [[ "$updated" == "0" ]]
then
  sed -i'.bak' '/^\s*gcAssistBytes/a \ \ \ \ \ \ \ \ routineCtx    uintptr' $RUNTIME_FILE
fi


if [ ! -f $TARGET_FILE ]
then
  echo "write code to $TARGET_FILE"
  sudo cat >  $TARGET_FILE  <<EOF
package runtime

func EllynGetGoid() uint64 {
  return uint64(getg().goid)
}

func EllynGetRoutineCtx() uintptr {
  return getg().routineCtx
}

func EllynSetRoutineCtx(ctx uintptr) {
  getg().routineCtx = ctx
}
EOF
  echo "write success"
else
  echo "goid file exists"
fi
