#!/bin/bash
set -e

GOMAXPROCS=1 go test -timeout 90s $(go list ./...)