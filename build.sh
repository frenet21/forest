#!/bin/bash
# Remove the databases
echo `rm -rf pubKeys`

# Remove the binary and rebuild
echo `rm -rf ./forest`
echo `go build -o forest`

echo "Build script completed."