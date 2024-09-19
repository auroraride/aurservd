#!/usr/bin/env bash

GOMODCACHE=$(go env GOMODCACHE)

dir() {
  echo "${GOMODCACHE}/$1@$(< go.mod sed $'/^require ($/,/^)$/!d; /^require ($/d;/^)$/d; /\\/\\/ indirect$/d; s/^\t+//g' | grep "$1" | cut -d' ' -f2)"
  return 0
}

ADAPTER="$(dir github.com/auroraride/adapter)"

APIS=(
  'controller/v1/common/common.go'
  'controller/v1/rapi/rapi.go'
  'controller/v2/rapi/rapi.go'
  'controller/v1/aapi/aapi.go'
  'controller/v1/eapi/eapi.go'
  'controller/v1/mapi/mapi.go'
  'controller/v2/mapi/mapi.go'
  'controller/v1/oapi/oapi.go'
  'controller/v1/papi/papi.go'
  'controller/v2/assetapi/assetapi.go'
  'controller/v2/wapi/wapi.go'
)

OUTPUTS=(
  './assets/docs/common/v1'
  './assets/docs/rider/v1'
  './assets/docs/rider/v2'
  './assets/docs/agent/v1'
  './assets/docs/employee/v1'
  './assets/docs/manager/v1'
  './assets/docs/manager/v2'
  './assets/docs/operator/v1'
  './assets/docs/promotion/v1'
  './assets/docs/asset/v2'
  './assets/docs/warestore/v2'
)

function exclude() {
  OUTPUT=''
  for API in "${APIS[@]}" ; do
    FOLDER_PATH=${API%/*}
    if [ "$1" != "$FOLDER_PATH" ]; then
      OUTPUT+=",./app/$FOLDER_PATH"
    fi
  done
  echo "$OUTPUT"
}

for i in "${!APIS[@]}" ; do
  API=${APIS[i]}
  FOLDER_PATH=${API%/*}
  TARGET_FILE=${API##*/}
  eval "swag fmt $TARGET_FILE"
  eval "swag init -g $TARGET_FILE -d ./app/$FOLDER_PATH,./app/model,./app/biz/definition,./app/permission,$ADAPTER -o ${OUTPUTS[i]} --md ./wiki"
done
