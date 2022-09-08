#!/bin/bash
cd internal/frontend
npm ci
npm run build
cd ../..
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o=bin/awp ./app.go