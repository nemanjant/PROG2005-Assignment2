FROM golang:1.19 as builder

LABEL maintainer="mail@domain.tld"
LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/cmd/server

# Copy relevant folders into container
COPY ./cmd /go/src/app/cmd
COPY ./data /go/src/app/data
COPY ./handler /go/src/app/handler
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# Indicate port on which server listens
EXPOSE 8080

# Instantiate binary
CMD ["./server"]




