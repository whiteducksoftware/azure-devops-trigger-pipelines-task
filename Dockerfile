# Builder
FROM golang:alpine as builder
WORKDIR /app

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Add src files
ADD . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/devops-worker

# Runner
FROM ubuntu:16.04

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apt update && \
    apt install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Add lables
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="Trigger Azure Pipelines" 
LABEL org.label-schema.description="Container which can trigger Azure pipeline(s)" 
LABEL org.label-schema.vcs-ref="https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task"
LABEL org.label-schema.maintainer="Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>"

# Copy our static executable.
COPY --from=builder /go/bin/devops-worker /go/bin/devops-worker
ENV PATH=$PATH:/go/bin