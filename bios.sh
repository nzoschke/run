#!/bin/bash
set -ex -o pipefail
env | grep -v PASS=

mkdir -p src/$PKG && cd src/$PKG && pwd

run -s "Cloning" git clone $URL --branch $REF --single-branch .
git branch --set-upstream-to=origin/$REF $REF
git reset --hard $SHA

PKGS=$(go list ./...)
run -s "Linting"  golint -set_exit_status $PKGS
run -s "Vetting"  go vet -x $PKGS
run -s "Building" go build -v $PKGS
run -s "Testing"  go test -v $PKGS
