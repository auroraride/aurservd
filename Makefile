define deploy
	swag init -g ./router/swagger.go -d ./app --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki --parseDependency --parseDepth 3
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/release/aurservd cmd/aurservd/main.go
	docker build --platform=linux/amd64 -t aurservd .
	docker tag aurservd registry.cn-beijing.aliyuncs.com/liasica/aurservd:$(1)
	docker push registry.cn-beijing.aliyuncs.com/liasica/aurservd:$(1)
endef

.PHONY: dev
dev:
	$(call deploy,latest)

.PHONY: prod
prod:
	$(call deploy,prod)