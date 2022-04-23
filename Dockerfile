FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/brunodeluk/kube-config/
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
COPY --from=builder /go/src/github.com/brunodeluk/kube-config/app ./
CMD ["./app"]
