FROM golang AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main ./cmd/kg/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app/kg/config/config.yaml ./app/kg/config/
COPY --from=builder /app/app/kg/config/fluentd.conf ./
COPY --from=builder /app/global/global.yaml ./global/
RUN mkdir -p /app/logs
RUN apk --no-cache add ca-certificates ruby ruby-dev build-base libc-dev linux-headers
RUN gem install fluentd --no-document
RUN gem install fluent-plugin-elasticsearch --no-document
RUN echo '#!/bin/sh' > ./start.sh && \
    echo 'fluentd -c ./fluentd.conf &' >> ./start.sh && \
    echo 'exec ./main' >> ./start.sh && \
    chmod +x ./start.sh
CMD ["./start.sh"]