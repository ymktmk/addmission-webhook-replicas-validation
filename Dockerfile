FROM --platform=linux/x86_64 golang:1.20 AS builder

WORKDIR /go/src/go-app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM --platform=linux/x86_64 alpine:latest

COPY --from=builder /go/src/go-app/main ./
