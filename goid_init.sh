#!/bin/sh

GOROOT=`go env GOROOT`
echo <<EOF
package runtime

func EllynGetGoid() uint64 {
	return uint64(getg().goid)
}

EOF > $GOROOT/ellyn_goid.go

