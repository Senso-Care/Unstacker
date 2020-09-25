#!/bin/bash

if [ "$#" -ne 3 ] || ! [ -d "./cmd/$1" ] || ! go tool dist list | grep "$2"/"$3" > /dev/null; then
  echo "Usage: $0 CMDNAME OS PLATFORM" >&2
  echo "Check os/platform available by running: \"go tool dist list\""
  exit 1
fi

cmd_path="./cmd/$1"

echo "Compiling $1 on $2/$3"
GOOS=$2 GOARCH=$3 go build -o="bin/$2/$3/$1" -mod=vendor "github.com/Senso-Care/daemons/$cmd_path"
