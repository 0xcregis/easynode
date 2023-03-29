FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags '-s -w' -o /app/easynode .


FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/cadok-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/easynode /app/easynode
COPY blockchain_config.json .
COPY collect_config.json .
COPY task_config.json .
COPY taskapi_config.json .

CMD ["./easynode","-collect_config","./collect_config.json","-task_config","./task_config.json","-blockchain_config","./blockchain_config.json","-taskapi_config","./taskapi_config.json"]
