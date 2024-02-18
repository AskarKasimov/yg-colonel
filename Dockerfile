FROM golang:alpine AS builder
WORKDIR /app
COPY / ./
RUN go build -o ./colonel ./main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/colonel ./colonel
COPY /templates ./templates
EXPOSE 8080
CMD ["./colonel"]