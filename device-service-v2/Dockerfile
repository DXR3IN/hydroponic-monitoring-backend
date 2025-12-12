FROM golang:1.24.3-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /device-service-v2 ./cmd/server
EXPOSE 8082
CMD ["/device-service-v2"]