FROM golang:1.22 as builder
ENV CGO_ENABLED=0

WORKDIR /build
COPY main.go go.* ./
RUN go mod download
RUN go build -o app

FROM alpine:latest
COPY --from=builder /build/app /usr/local/bin/app

CMD ["app"]