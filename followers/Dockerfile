FROM golang:alpine AS followers-builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o followers-webapp

FROM alpine
COPY --from=tours-builder /app/followers-webapp /usr/bin/followers-webapp
EXPOSE 8080
ENTRYPOINT ["/usr/bin/followers-webapp"]