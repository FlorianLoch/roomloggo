# The Raspberry will require an arm image, so the platform building the image should be arm too. If this
# is not possible, additioniol measures will probably be required - like building with `docker buildx`
# Using an alpine based image for the build is not possible, we need the full toolchain because of cgo
FROM golang:1.21.4 AS builder

WORKDIR /src/
# Copy dependency management related files first and download required modules before copying changed code into the
# container. That way we can cache the downloading as long as the dependency configuration does not change too.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# We need a static binary as CGO_ENABLED=1 and dynamic linking does not work using alpine as runtime environment
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o roomloggo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /src/roomloggo ./
CMD ["./roomloggo"]
