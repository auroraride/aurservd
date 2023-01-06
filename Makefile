# GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/release/aurservd cmd/aurservd/main.go
# GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/release/aurservd cmd/aurservd/main.go

define deploy
	swag init -g ./router/docs.go -d ./app --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki --parseDependency --parseDepth 100
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o build/release/aurservd cmd/aurservd/main.go
	docker build --platform=linux/amd64 -t aurservd .
	docker tag aurservd registry.cn-beijing.aliyuncs.com/liasica/aurservd:$(1)
	docker push registry.cn-beijing.aliyuncs.com/liasica/aurservd:$(1)
endef

.PHONY: dev
dev:
	$(call deploy,dev)

.PHONY: prod
prod:
	$(call deploy,prod)


.PHONY: latest
latest:
	$(call deploy,latest)


.PHONY: permisson
permission:
	go run ./cmd/permission