#!/bin/bash

set -a
[ -f .env ] && source .env
set +a

case "$1" in
  api)
    cd go
    go run cmd/api/main.go
    ;;
  migrate-up)
    cd go
    go run cmd/migrations/*.go
    ;;
  migrate-init)
    cd go
    go run cmd/migrations/main.go init
    ;;
  migrate-version)
    cd go
    go run cmd/migrations/main.go version
    ;;
  docker-up)
    docker-compose up -d
    ;;
  docker-down)
    docker-compose down
    ;;
  docker-build)
    docker build --progress=plain --no-cache -f Postgres.dockerfile -t postgres-custom .
    ;;
  *)
    echo "Usage: $0 {docker-up|docker-down|docker-build|migrate-up|migrate-init}"
    exit 1
    ;;
esac 