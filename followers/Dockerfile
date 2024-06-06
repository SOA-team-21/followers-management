FROM golang:alpine AS followers-builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o followers-webapp

FROM alpine
COPY --from=followers-builder /app/followers-webapp /usr/bin/followers-webapp
EXPOSE 87
ENTRYPOINT ["/usr/bin/followers-webapp"]