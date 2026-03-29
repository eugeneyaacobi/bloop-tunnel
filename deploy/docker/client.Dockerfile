FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/bloop-tunnel ./cmd/bloop-tunnel

FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/bloop-tunnel /bloop-tunnel
ENTRYPOINT ["/bloop-tunnel"]
