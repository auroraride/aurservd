FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && mkdir /app \
    && apk add --no-cache bash tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*
COPY ./build/release/aurservd /app/

WORKDIR /app

ENTRYPOINT ["/app/aurservd", "server"]