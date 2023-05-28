FROM docker.io/golang:1.20-bullseye as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
RUN git clone --depth 1 --single-branch --branch main https://github.com/babinkos/webhook-bb-pr-teams-router-app-go.git
RUN cd webhook-bb-pr-teams-router-app-go/backend-test && go mod download && go mod verify

# COPY main.go ./
RUN cd webhook-bb-pr-teams-router-app-go/backend-test && go build -v -o /usr/local/bin/app ./...

FROM docker.io/debian:bullseye-slim
COPY --from=build /usr/local/bin/app /usr/local/bin/app
EXPOSE 3000
CMD ["/usr/local/bin/app"]