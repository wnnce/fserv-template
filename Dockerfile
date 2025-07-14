# syntax=docker/dockerfile:1
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ----------------------------
# Final image with timezone & fonts
# ----------------------------
FROM alpine:latest

RUN apk --no-cache add \
    tzdata \
    ca-certificates \
    font-noto \
 && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
