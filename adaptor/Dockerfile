FROM docker.io/golang:1.20-bookworm as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
RUN git clone --depth 1 --single-branch --branch main https://github.com/babinkos/webhook-bb-pr-teams-router-app-go.git
RUN \
    cd webhook-bb-pr-teams-router-app-go/adaptor && \
    go mod download && go mod verify

# COPY main.go ./
ENV DEBIAN_FRONTEND=noninteractive
RUN \
    apt update -y && \
    apt upgrade -y

RUN \
    cd webhook-bb-pr-teams-router-app-go/adaptor && \
    GOOS=linux go build -tags netgo -ldflags "-w -s -linkmode external -extldflags -static" -v -o /usr/local/bin/app main.go

FROM scratch
COPY --from=build /usr/local/bin/app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
EXPOSE 9000
USER 10000
ENTRYPOINT ["/app"]
LABEL org.opencontainers.image.description="Distroless image with adaptor" \
  org.opencontainers.image.documentation="https://github.com/babinkos/webhook-bb-pr-teams-router-app-go/blob/main/README.md" \
  org.opencontainers.image.base.name="ghcr.io/babinkos/webhook-bb-pr-teams-router-app-go:main" \
  org.opencontainers.image.url="https://github.com/babinkos/webhook-bb-pr-teams-router-app-go/blob/main/adaptor/Dockerfile"
