# --- Build Stage ---
FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .
RUN go build -o server ./cmd

# --- Production Stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/server .

EXPOSE 8000
CMD ["./server"]
