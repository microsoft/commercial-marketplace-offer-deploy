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
COPY ./build/scripts/entrypoint.sh ./bin


# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
#FROM debian:buster-slim
FROM registry.access.redhat.com/ubi9/nginx-120
#RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
#    ca-certificates && \
#    rm -rf /var/lib/apt/lists/*
USER root

ARG ACME_RELEASE_TAG=3.0.5
ARG DRIVER_RELEASE_TAG=dev

RUN yum -y --repo ubi-9-appstream-rpms install socat && \
  curl -L -o acme.zip https://github.com/acmesh-official/acme.sh/archive/refs/tags/v${ACME_RELEASE_TAG}.zip && \
  unzip -qoj acme.zip acme.sh-${ACME_RELEASE_TAG}/acme.sh -d . && rm acme.zip && \
  echo "ACME=${ACME_RELEASE_TAG}" >> versions && echo "DRIVER=${DRIVER_RELEASE_TAG}" >> versions && ls -al && cat versions

ENV LANG en_US.utf8


ADD ["/templates/nginx", "/etc/nginx/"]

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/operator /operator
COPY --from=builder /app/bin/apiserver /apiserver
COPY --from=builder /app/bin/entrypoint.sh /entrypoint.sh

RUN chmod +x /operator
RUN chmod +x /apiserver
RUN chmod +x /entrypoint.sh

# Run the web service on container startup.
ENTRYPOINT ["/entrypoint.sh"]
