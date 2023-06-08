FROM golang:1.20.4

USER root

# Delve + Hot reloading
RUN go install github.com/go-delve/delve/cmd/dlv@v1.8.1
RUN go install github.com/cosmtrek/air@v1.29.0
ENV GOPATH=/root/go
ENV PATH="${PATH}:/root/go/bin"

# Application
COPY . /opt/app-root/
WORKDIR /opt/app-root/

# Cache dependencies in Docker image
RUN go mod download

CMD air -c .air.toml
