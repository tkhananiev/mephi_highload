FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:3.20
WORKDIR /app

COPY --from=build /src/app /app/go-microservice

EXPOSE 8080
ENV PORT=8080

CMD ["./go-microservice"]
