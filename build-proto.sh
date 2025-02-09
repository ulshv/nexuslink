#!/bin/bash

# Recursively find and build all .proto files from the current directory
find . -name "*.proto" -type f -exec protoc -I=. --go_out=. {} \;
