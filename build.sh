#!/bin/bash
# Remove the databases
echo "Deleting databases..."
echo `rm -rf .pubKeys .priKeys .knownBlocks`

# Remove the binary and rebuild
echo "Deleting existing binaries..."
echo `rm -rf ./forest`

echo "Getting external packages..."
echo `go get golang.org/x/crypto/sha3`
echo `go get github.com/skratchdot/open-golang/open`
echo `go get github.com/syndtr/goleveldb/leveldb`

echo "Building forest..."
echo `go build -o forest`

echo "Build script completed."