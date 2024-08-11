#!/bin/bash

clear && go run main.go

PID=$!

while inotifywait -e modify -r ./*.go; do
  echo "Process Reset"
  kill $PID

  clear && go run main.go

  PID=$!

done
