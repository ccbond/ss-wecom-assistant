FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ss-wecom-assistant .

FROM alpine:latest  
RUN apk update && apk --no-cache add ca-certificates
COPY --from=builder /app/ss-wecom-assistant /ss-wecom-assistant
COPY ./conf/config.toml /data/ss-wecom-assistant/config.toml
COPY ./asset/qwewm.jpg /data/ss-wecom-assistant/qwewm.jpg

EXPOSE 1107
ENTRYPOINT ["/ss-wecom-assistant", "-e", "online"]
