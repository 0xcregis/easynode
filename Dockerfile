FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags '-s -w' -o /app/easynode ./cmd/easynode/app.go


FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/cadok-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


WORKDIR /app
CMD  mkdir config
COPY --from=builder /app/easynode /app/easynode

EXPOSE 9001 9002 9003

ENTRYPOINT ["./easynode","-collect","./config/collect_config.json","-task","./config/task_config.json","-blockchain","./config/blockchain_config.json","-taskapi","./config/taskapi_config.json","-store","./config/store_config.json"]
