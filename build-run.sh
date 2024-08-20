#!/bin/bash

echo "Build vue proyect"
npm run --prefix ui/frontend/ build

echo "Remove previous app build"
rm web

echo "Build app and run"
go build ./cmd/web/
./web
