#!/bin/bash
# Remove the databases
echo "Deleting databases..."
echo `rm -rf pubKeys`

# Remove the binary and rebuild
echo "Deleting binaries..."
echo `rm -rf ./forest`

echo "Building forest..."
echo `go build -o forest`

echo "Build script completed."