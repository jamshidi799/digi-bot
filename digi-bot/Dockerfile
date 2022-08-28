FROM golang:1.18-alpine AS build_base

RUN apk add build-base
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/go-sample-app -tags musl .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/out/go-sample-app /app/go-sample-app

# Run the binary program produced by `go install`
CMD ["/app/go-sample-app"]