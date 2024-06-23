FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ss-wecom-assistant .

FROM scratch
COPY --from=builder /app/ss-wecom-assistant /ss-wecom-assistant
COPY ./conf/config.toml /usr/conf/config.toml
COPY ./asset/qwewm.jpg /usr/conf/qwewm.jpg

EXPOSE 1107
ENTRYPOINT ["/ss-wecom-assistant", "-e", "online"]
