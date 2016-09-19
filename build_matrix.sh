#!/bin/bash

for CLIAPP in we toconfig bach; do
  for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
      echo "Building $GOOS-$GOARCH"
      export GOOS=$GOOS
      export GOARCH=$GOARCH
      export CLIAPP=$CLIAPP
      make bin/$CLIAPP-$GOOS-$GOARCH
    done
  done
done
