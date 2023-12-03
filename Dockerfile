#Build Stage


FROM golang:1.21 AS builder
WORKDIR /ytb
COPY . .
RUN go mod download
EXPOSE 8000
RUN go build .

#Final Stage
FROM alpine:latest
RUN apk add --no-cache \
  ffmpeg \
  tzdata

WORKDIR /
COPY --from=builder . .
EXPOSE 8000

CMD ["./ytb"]