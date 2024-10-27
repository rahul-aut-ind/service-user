# ---- Build ----
FROM golang:1.22-alpine AS golang-base
ENV SERVICE_NAME="service-user"

WORKDIR /go/src/${SERVICE_NAME}
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=arm64/v8 cd ./cmd/${SERVICE_NAME} && go build -a -o ../../dist/service

# ---- Final ----
FROM alpine:latest
ENV SERVICE=service
ENV SERVER_PORT=8080

EXPOSE ${SERVER_PORT}
WORKDIR /app
COPY --from=golang-base /go/src/service-user/dist .
ENTRYPOINT []
CMD "./${SERVICE}"