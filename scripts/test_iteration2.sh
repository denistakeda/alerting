#!/bin/bash

./scripts/buildall.sh

devopstest -test.v -test.run=^TestIteration2[b]*$ \
  -source-path=. \
  -binary-path=cmd/server/server
