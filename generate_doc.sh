#!/usr/bin/env bash

GOMODCACHE=$(go env GOMODCACHE)

dir() {
  echo "${GOMODCACHE}/$1@$(< go.mod sed $'/^require ($/,/^)$/!d; /^require ($/d;/^)$/d; /\\/\\/ indirect$/d; s/^\t+//g' | grep "$1" | cut -d' ' -f2)"
  return 0
}

ADAPTER="$(dir github.com/auroraride/adapter)"

swag init -g controller/v1/rapi/rapi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/aapi,./app/controller/v1/eapi,./app/controller/v1/mapi,./app/controller/v1/oapi,./app/controller/v1/papi -o ./assets/docs/rider/v1 --md ./wiki
swag init -g controller/v1/aapi/aapi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/rapi,./app/controller/v1/eapi,./app/controller/v1/mapi,./app/controller/v1/oapi,./app/controller/v1/papi -o ./assets/docs/agent/v1 --md ./wiki
swag init -g controller/v1/eapi/eapi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/aapi,./app/controller/v1/rapi,./app/controller/v1/mapi,./app/controller/v1/oapi,./app/controller/v1/papi -o ./assets/docs/employee/v1 --md ./wiki
swag init -g controller/v1/mapi/mapi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/aapi,./app/controller/v1/eapi,./app/controller/v1/rapi,./app/controller/v1/oapi,./app/controller/v1/papi -o ./assets/docs/manager/v1 --md ./wiki
swag init -g controller/v1/oapi/oapi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/aapi,./app/controller/v1/eapi,./app/controller/v1/mapi,./app/controller/v1/rapi,./app/controller/v1/papi -o ./assets/docs/operator/v1 --md ./wiki
swag init -g controller/v1/papi/papi.go -d ./app,"$ADAPTER" --exclude ./app/service,./app/router,./app/middleware,./app/request,./app/controller/v1/aapi,./app/controller/v1/eapi,./app/controller/v1/mapi,./app/controller/v1/oapi,./app/controller/v1/rapi -o ./assets/docs/promotion/v1 --md ./wiki
