FROM alpine:latest

RUN mkdir /app
COPY ./build/release/aurservd /app/

WORKDIR /app

ENTRYPOINT ["/app/aurservd"]