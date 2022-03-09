FROM alpine:latest

RUN mkdir /app && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/localtime && apk del tzdata
COPY ./build/release/aurservd /app/

WORKDIR /app

ENTRYPOINT ["/app/aurservd", "server"]