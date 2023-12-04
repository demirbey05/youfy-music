#Build Stage
FROM golang:1.21 AS builder
WORKDIR /ytb
COPY . .
RUN go mod download
EXPOSE 80
RUN go build .

#Final Stage
FROM alpine:latest
RUN set -x \
 && apk add --no-cache ca-certificates curl ffmpeg python3 \
    # Install youtube-dlp
 && curl -Lo /usr/local/bin/youtube-dl https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
 && chmod a+rx /usr/local/bin/youtube-dl \
    # Clean-up
 && apk del curl \
    # Create directory to hold downloads.
 && mkdir /downloads \
 && chmod a+rw /downloads \
    # Basic check it works.
 && youtube-dl --version \
 && mkdir -p /.cache \
 && chmod 777 /.cache

WORKDIR /
COPY --from=builder . .
EXPOSE 80

CMD ["./ytb"]