# Builder
FROM golang:alpine as builder
WORKDIR /app

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Add src files
ADD . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/devops-worker

# Runner
FROM scratch

# Add lables
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="Trigger Azure Pipelines" 
LABEL org.label-schema.description="Container which can trigger Azure pipeline(s)" 
LABEL org.label-schema.vcs-ref="https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task"
LABEL org.label-schema.maintainer="Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>"

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/devops-worker /go/bin/devops-worker
ENV PATH=$PATH:/go/bin

# Use an unprivileged user.
USER appuser:appuser