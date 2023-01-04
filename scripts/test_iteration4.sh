#!/bin/bash

./scripts/buildall.sh

devopstest -test.v -test.run=^TestIteration4$ \
    -source-path=. \
    -binary-path=cmd/server/server \
    -agent-binary-path=cmd/agent/agent