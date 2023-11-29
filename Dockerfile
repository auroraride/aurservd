FROM alpine
RUN  sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
     && mkdir /app \
     && apk add --no-cache bash tzdata \
     && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && echo "Asia/Shanghai" > /etc/timezone \
     && rm -rf /var/cache/apk/* \
     && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY build/release/aurservd /app/aurservd
ENTRYPOINT ["/app/aurservd", "server"]
