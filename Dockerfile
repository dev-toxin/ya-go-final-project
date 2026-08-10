FROM golang:1.24-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler ./main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=build /scheduler /scheduler
COPY web ./web

ENTRYPOINT ["/scheduler"]
