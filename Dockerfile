FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN cp /build/.env .

FROM scratch

COPY --from=builder /dist/main .
COPY --from=builder /dist/.env .

ENTRYPOINT ["/main"]