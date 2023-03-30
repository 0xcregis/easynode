FROM golang:1.18 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download
COPY . .
RUN go build -ldflags '-s -w' -o /app/easynode ./cmd/easynode/app.go


FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/cadok-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/easynode /app/easynode
COPY ./cmd/easynode/blockchain_config.json .
COPY ./cmd/easynode/collect_config.json .
COPY ./cmd/easynode/task_config.json .
COPY ./cmd/easynode/taskapi_config.json .

EXPOSE 9001 9002

CMD ["./easynode","-collect_config","./collect_config.json","-task_config","./task_config.json","-blockchain_config","./blockchain_config.json","-taskapi_config","./taskapi_config.json"]
