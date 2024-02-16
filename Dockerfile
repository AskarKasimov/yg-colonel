FROM golang:alpine AS builder
WORKDIR /app
COPY / ./
RUN apk add build-base && apk cache clean
ENV CGO_ENABLED=1
RUN go build -o ./colonel ./main.go


FROM alpine
WORKDIR /app
COPY --from=builder /app/colonel ./colonel
EXPOSE 8080
CMD ["./colonel"]