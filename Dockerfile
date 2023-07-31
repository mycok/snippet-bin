FROM golang:1.20-alpine3.17 AS build
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build -a \
    ldflags "-extldflags '-static' -w -s" \
    -o /bin/app /src/

FROM alpine3.17
COPY --from=build /bin/app /bin/app
EXPOSE 4000
CMD [ "/bin/app" ]