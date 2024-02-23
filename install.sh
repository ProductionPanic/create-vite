#!/bin/bash
# check if npm is installed
if ! [ -x "$(command -v npm)" ]; then
  echo 'Error: npm is not installed.' >&2
  exit 1
fi

# check if go is installed
if ! [ -x "$(command -v go)" ]; then
  echo 'Error: go is not installed.' >&2
  exit 1
fi

go build .

sudo mv ./create-vite /usr/local/bin/create-vite

echo "create-vite installed successfully"
