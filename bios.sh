#!/bin/bash
set -ex -o pipefail

export GIT_DIR=src/$PKG/.git
run -s "Cloning"      git clone $URL --branch $REF --single-branch src/$PKG
run -s "Resetting"    git reset --hard $SHA
run -s "Fetching"     git fetch origin $BREF
run -s "Whitespacing" git diff-tree --check $BSHA $SHA

PKGS=$(go list $PKG/...)
run -s "Linting"  golint -set_exit_status $PKGS
run -s "Vetting"  go vet -x $PKGS
run -s "Building" go build -v $PKGS
run -s "Testing"  go test -v $PKGS
