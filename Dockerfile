FROM alpine:latest

RUN apk add --no-cache ffmpeg

ENV GOLANG_VERSION 1.20
RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go
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
RUN CGO_ENABLED=0 GOOS=linux go build -o /ss-assistant

EXPOSE 80
ENTRYPOINT ["/ss-assistant", "-e", "online"]
