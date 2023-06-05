#!/usr/bin/env bash

GOMODCACHE=$(go env GOMODCACHE)

dir() {
  echo "${GOMODCACHE}/$1@$(< go.mod sed $'/^require ($/,/^)$/!d; /^require ($/d;/^)$/d; /\\/\\/ indirect$/d; s/^\t+//g' | grep $1 | cut -d' ' -f2)"
  return 0
}

ADAPTER="$(dir github.com/auroraride/adapter)"

swag init -g ./router/docs.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki
