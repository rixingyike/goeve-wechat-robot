#!/bin/bash
# Platforms: android, darwin, ios, linux, windows
# Achitectures: 386, amd64, arm-5, arm-6, arm-7, arm64, mips, mipsle, mips64, mips64le

xgo --targets=darwin/amd64 -out ./release/sim-robot .

# go build -o ./release/sim-robot_darwin main.go
echo "complete."