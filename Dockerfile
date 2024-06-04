FROM alpine:latest

ENV GOLANG_VERSION 1.21.1
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
RUN cd /usr && mkdir conf && cd ..

COPY . .
COPY ./conf/config.toml /usr/conf/config.toml
COPY ./asset/qwewm.jpg /usr/conf/qwewm.jpg

RUN CGO_ENABLED=0 GOOS=linux go build -o /ss-wecom-assistant

EXPOSE 1107
CMD ["/ss-wecom-assistant", "-e", "online"]
