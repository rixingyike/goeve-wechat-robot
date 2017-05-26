#!/bin/bash
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=exectuable-of-mingw-w64-gcc go build -o ./release/sim-robot.exe main.go
echo "complete."