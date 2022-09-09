#!/bin/bash
echo "Building version $1"
git tag $1
cd internal/frontend
npm ci
npm run build
cd ../..
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/awp ./app.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=$1'" -o=bin/app_darwin_amd64 ./app.go