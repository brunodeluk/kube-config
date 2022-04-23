FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/brunodeluk/kube-config/
COPY ./ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
COPY --from=builder /go/src/github.com/brunodeluk/kube-config/app ./
CMD ["./app"]
