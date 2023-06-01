FROM --platform=linux/amd64 golang as builder

WORKDIR /usr/src/app

COPY . .

RUN git clone https://github.com/liasica/swag.git && cd swag/cmd/swag && go install .
RUN mkdir -p ./assets/docs && \
    go get ./... && \
    swag init -g ./router/docs.go -d ./app --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki --parseDependency && \
    CGO_ENABLED=0 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o /go/bin/aurservd cmd/aurservd/main.go


FROM scratch
COPY --from=builder /go/bin/aurservd /app/aurservd
COPY --from=builder /usr/share/zoneinfo /usr/share/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT ["/app/aurservd", "server"]