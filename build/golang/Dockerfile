FROM golang:1.20 AS build-server

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /output/tanker-server /app/cmd/server/server.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /output/tanker-cli /app/cmd/cli/cli.go

FROM alpine:latest

WORKDIR /app

COPY --from=build-server /app/configs/config.yaml /app/config.yaml
COPY --from=build-server /output/tanker-server /app/tanker-server
COPY --from=build-server /output/tanker-cli /app/tanker-cli

RUN apk add tzdata

RUN chmod 755 /app/tanker-server
RUN chmod +x /app/tanker-server
RUN chmod +x /app/tanker-cli

EXPOSE 80

ENTRYPOINT ["/app/tanker-server", "/app/config.yaml"]