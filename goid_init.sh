#!/bin/sh

GOROOT=`go env GOROOT`

TARGET_FILE=$GOROOT/src/runtime/ellyn_goid.go

echo "write code to $TARGET_FILE"

echo <<EOF
package runtime

func EllynGetGoid() uint64 {
	return uint64(getg().goid)
}

EOF > $TARGET_FILE

echo "write success"