# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

# Copy local code to the container image.
COPY . ./

RUN go build -o ./bin/ ./cmd/operator
RUN go build -o ./bin/ ./cmd/apiserver
COPY ./scripts/startservices.sh ./bin

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV LANG en_US.utf8

# install blobfuse
RUN apt-get update \
    && apt-get install -y wget apt-utils \
    && wget https://packages.microsoft.com/config/ubuntu/18.04/packages-microsoft-prod.deb \
    && dpkg -i packages-microsoft-prod.deb \
    && apt-get remove -y wget \
    && apt-get update \
    && apt-get install -y --no-install-recommends fuse blobfuse libcurl3-gnutls libgnutls30 \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/operator /operator
COPY --from=builder /app/bin/apiserver /apiserver
COPY --from=builder /app/bin/startservices.sh /startservices.sh

RUN chmod +x ./startservices.sh
RUN chmod +x ./operator
RUN chmod +x ./apiserver

# Run the web service on container startup.
ENTRYPOINT ["./startservices.sh"]
