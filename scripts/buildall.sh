#!/bin/bash

rm cmd/server/server
cd cmd/server
go build -o server
cd ../..

rm cmd/agent/agent
cd cmd/agent
go build -o agent
cd ../..