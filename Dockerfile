#Builder stage
FROM golang:alpine AS builder

#Install dependencies
RUN apk update && apk add --no-cache \
  git 
# ca-certificates \
# && update-ca-certificates

# Add source files and set the proper work dir
COPY . $GOPATH/src/github.com/bysidecar/untraceable/
WORKDIR $GOPATH/src/github.com/bysidecar/untraceable/cmd


# Enable Go modules
ENV GO111MODULE=on
# Build the binary
RUN go build -mod=vendor -o /go/bin/livelead

# Final image
FROM alpine

# Copy our static executable
COPY --from=builder /go/bin/untraceable /go/bin/untraceable

# Copy the ca-certificates to be able to perform https requests
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the binary
ENTRYPOINT ["/go/bin/untraceable"]