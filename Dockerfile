FROM --platform=linux/amd64 golang:alpine as builder

WORKDIR /usr/src/app

COPY . .

#RUN wget -O - -q https://github.com/liasica/swag/releases/download/v1.16.1-733b9ee/swag_1.16.1-733b9ee_Linux_x86_64.tar.gz | tar -xz -C /go/bin/
#RUN git clone https://github.com/liasica/swag.git && cd swag/cmd/swag && go install .
RUN apk --no-cache add tzdata git
RUN wget -O - -q https://github.com/liasica/swag/releases/download/v1.16.1-733b9ee/swag_1.16.1-733b9ee_Linux_x86_64.tar.gz | tar -xz -C /go/bin/
RUN go get ./... && \
    bash ./generate_doc.sh && \
    CGO_ENABLED=0 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o /go/bin/aurservd cmd/aurservd/main.go


FROM alpine
COPY --from=builder /go/bin/aurservd /app/aurservd
#COPY --from=builder /usr/share/zoneinfo /usr/share/
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN  apk add --no-cache bash tzdata \
     && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && echo "Asia/Shanghai" > /etc/timezone \
     && rm -rf /var/cache/apk/* \
     && rm -rf /var/lib/apt/lists/*
WORKDIR /app
ENTRYPOINT ["/app/aurservd", "server"]
