#!/bin/bash

# Load test environment variables
set -a
source .env.test
set +a

# Run unit tests with coverage
go test -v -count=1 -cover ./internal/... ./pkg/... ./models/... 