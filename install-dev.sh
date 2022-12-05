#!/bin/bash
version="DEV_EDITION"
echo "Building developer version"
cd internal/proxy/frontend
npm ci
npm run build
cd ../../..
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/tfmcdigital/aws-web-proxy/internal/domain.Version=$version'" -o=bin/awp ./app.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/tfmcdigital/aws-web-proxy/internal/domain.Version=$version'" -o=bin/app_darwin_amd64 ./app.go
cp bin/awp ~/awp/awp