#!/bin/bash
env GOOS=windows GOARCH=amd64 go build -o ./release/sim-robot_win64.exe main.go
echo "complete."