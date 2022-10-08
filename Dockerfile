FROM golang:1.19.1 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o NFT-Sales-Bot .


FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/NFT-Sales-Bot .

ENTRYPOINT ["./NFT-Sales-Bot"]