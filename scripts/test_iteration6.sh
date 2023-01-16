#!/bin/bash

./scripts/buildall.sh

SERVER_PORT=3000
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=/tmp/devops-metrics-db.json
devopstest -test.v -test.run=^TestIteration6$ \
  -source-path=. \
  -agent-binary-path=cmd/agent/agent \
  -binary-path=cmd/server/server \
  -server-port=$SERVER_PORT \
  -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
  -file-storage-path=$TEMP_FILE
