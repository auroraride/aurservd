FROM --platform=linux/amd64 golang as builder

WORKDIR /usr/src/app

COPY . .

#RUN wget -O - -q https://github.com/liasica/swag/releases/download/v1.16.1-733b9ee/swag_1.16.1-733b9ee_Linux_x86_64.tar.gz | tar -xz -C /go/bin/
RUN git clone https://github.com/liasica/swag.git && cd swag/cmd/swag && go install .
RUN go get ./... && \
    swag init -g ./router/docs.go -d ./app --exclude ./app/service,./app/router,./app/middleware,./app/request -o ./assets/docs --md ./wiki --parseDependency && \
    CGO_ENABLED=0 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o /go/bin/aurservd cmd/aurservd/main.go


FROM scratch
COPY --from=builder /go/bin/aurservd /app/aurservd
COPY --from=builder /usr/share/zoneinfo /usr/share/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT ["/app/aurservd", "server"]