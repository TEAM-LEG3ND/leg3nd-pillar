FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates


COPY . .

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN cp /build/.env .

FROM scratch

COPY --from=builder /dist/main .
COPY --from=builder /dist/.env .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENTRYPOINT ["/main"]