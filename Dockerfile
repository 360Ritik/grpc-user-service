# Use golang as builder stage
FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-user-service .

# Use distroless as minimal base image
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/grpc-user-service /app/grpc-user-service

EXPOSE 50051

CMD ["/app/grpc-user-service"]
