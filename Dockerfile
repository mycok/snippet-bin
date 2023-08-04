FROM golang:1.20-alpine3.17 AS builder
WORKDIR /src/
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build -a \
    -ldflags "-extldflags '-static' -w -s" \
    -o /src/bin/app ./cmd/web

FROM alpine:3.17
COPY --from=builder /src/bin/app /src/bin/app
EXPOSE 4000
CMD [ "/src/bin/app" ]