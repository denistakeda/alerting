#!/bin/bash

./scripts/buildall.sh

SERVER_PORT=3000
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=/tmp/devops-metrics-db.json
devopstest -test.v -test.run=^TestIteration9$ \
  -source-path=. \
  -agent-binary-path=cmd/agent/agent \
  -binary-path=cmd/server/server \
  -server-port=$SERVER_PORT \
  -file-storage-path=$TEMP_FILE \
  -database-dsn='***postgres:5432/praktikum?sslmode=disable' \
  -key="${TEMP_FILE}"