FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ss-wecom-assistant .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 1107
RUN cd /usr && mkdir conf
COPY ./conf/config.toml /usr/conf/config.toml
COPY --from=builder /app/ss-wecom-assistant .
CMD ["./ss-wecom-assistant", "-e", "online"]
