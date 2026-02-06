# Build UI
FROM node:18 as ui-builder
WORKDIR /app/web
COPY web/package.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Build Server
FROM golang:1.21 as server-builder
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

# Final Image
FROM alpine:latest
WORKDIR /root/
COPY --from=server-builder /app/server .
COPY --from=ui-builder /app/web/build ./web/build
EXPOSE 8080 9090
CMD ["./server"]
