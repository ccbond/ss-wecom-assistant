FROM alpine:latest  

ENV GOLANG_VERSION 1.20
RUN apk update && apk --no-cache --virtual .build-deps bash gcc musl-dev openssl go ca-certificates
RUN wget -O go.tgz "https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz" && \
    tar -C /usr/local -xzf go.tgz && \
    cd /usr/local/go/src && \
    ./make.bash && \
    rm -rf /go.tgz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ss-wecom-assistant

EXPOSE 1107
ENTRYPOINT ["/ss-wecom-assistant", "-e", "online"]
