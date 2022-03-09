FROM alpine:latest

RUN mkdir /app \
    && apk add --no-cache bash tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/localtime \
    && apk del tzdata \
    && rm -rf /var/cache/apk/*
COPY ./build/release/aurservd /app/

WORKDIR /app

ENTRYPOINT ["/app/aurservd", "server"]