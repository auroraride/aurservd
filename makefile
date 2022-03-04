.PHONY: build
build:
	swag init -g ./router/swagger.go -d ./app --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/release/aurservd cmd/aurservd/main.go
	docker build -t aurservd .
	docker tag aurservd registry.cn-beijing.aliyuncs.com/liasica/aurservd:latest
	docker push registry.cn-beijing.aliyuncs.com/liasica/aurservd:latest