FROM golang:1.19-alpine as builder

WORKDIR /app

# copy module manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache modules
RUN go mod download

# copy source code
COPY main.go main.go

# build
RUN CGO_ENABLED=0 go build \
    -a -o leaderelection main.go

FROM alpine:3.17.1

RUN apk --no-cache add ca-certificates

USER nobody

COPY --from=builder --chown=nobody:nobody /app/leaderelection .

ENTRYPOINT ["./leaderelection"]