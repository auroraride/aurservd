.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/release/aurservd cmd/aurservd/main.go
	docker build -t aurservd .
	docker tag aurservd registry.cn-beijing.aliyuncs.com/liasica/aurservd:latest
	docker push registry.cn-beijing.aliyuncs.com/liasica/aurservd:latest