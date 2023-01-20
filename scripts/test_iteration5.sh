#!/bin/bash

./scripts/buildall.sh

SERVER_PORT=3000
ADDRESS="localhost:${SERVER_PORT}"
devopstest -test.v -test.run=^TestIteration5$ \
    -source-path=. \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -server-port=$SERVER_PORT