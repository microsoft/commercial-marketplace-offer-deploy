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

#RUN go build -o ./bin/ ./tools/testharness/main.go
RUN go build -o ./bin/ ./tools/testharness
COPY ./scripts/start-testharness.sh ./bin/
COPY ./test/testdata/taggeddeployment/mainTemplateBicep.json ./bin
COPY ./test/testdata/taggeddeployment/parametersBicep.json ./bin

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM registry.access.redhat.com/ubi9/nginx-120

USER root

# RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
#    ca-certificates && \
#    rm -rf /var/lib/apt/lists/*

ENV LANG en_US.utf8

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/testharness /testharness
COPY --from=builder /app/bin/mainTemplateBicep.json /mainTemplateBicep.json
COPY --from=builder /app/bin/parametersBicep.json /parametersBicep.json
COPY --from=builder /app/bin/start-testharness.sh /start-testharness.sh

RUN chmod +x /testharness
RUN chmod +x /start-testharness.sh

# Run the web service on container startup.
ENTRYPOINT ["/start-testharness.sh"]
