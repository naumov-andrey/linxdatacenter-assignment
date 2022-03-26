FROM golang:1.17-alpine AS builder

WORKDIR /workdir

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o build/main cmd/main.go


FROM alpine:latest

WORKDIR /workdir

COPY --from=builder /workdir/build/main ./
ENTRYPOINT ["./main"]
