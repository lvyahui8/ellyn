#!/bin/bash

set -e

GOROOT=`go env GOROOT`

TARGET_FILE=$GOROOT/src/runtime/ellyn_goid.go

if [ ! -f $TARGET_FILE ]
then
  echo "write code to $TARGET_FILE"
  sudo cat >  $TARGET_FILE  <<EOF
  package runtime

  func EllynGetGoid() uint64 {
    return uint64(getg().goid)
  }
EOF
  echo "write success"
else
  echo "goid file exists"
fi
