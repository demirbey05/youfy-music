# Build Stage
FROM golang:1.20 AS builder
WORKDIR /ytb
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ytb .

# Final Stage
FROM alpine:latest
RUN set -x \
 && apk add --no-cache ca-certificates curl ffmpeg python3 \
    # Install youtube-dlp
 && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
 && chmod a+rx /usr/local/bin/yt-dlp \
    # Clean-up
 && apk del curl \
    # Create directory to hold downloads.
 && yt-dlp --version \


WORKDIR /
COPY --from=builder /ytb/ .
EXPOSE 80

CMD ["./ytb"]
