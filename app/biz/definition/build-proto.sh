#!/usr/bin/env sh

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative "$(find . -iname '*.proto')" --dart_out=/Users/liasica/projects/auroraride/v2/rider/lib/modules/face/data/gen
