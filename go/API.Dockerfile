FROM golang:1.24.4 as builder

WORKDIR /app

# Copy the local package files to the container's workspace.
ADD ./go ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add curl
COPY --from=builder /app/server ./server

CMD ["./server"]

EXPOSE 80