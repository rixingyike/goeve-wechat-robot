#!/bin/bash
env GOOS=windows GOARCH=386 go build -o ./release/sim-robot_win32.exe main.go
echo "complete."