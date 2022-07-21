# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS build
WORKDIR /
COPY ./* ./
RUN go build -o /app
RUN chmod +x /app

##

FROM alpine:3.14
WORKDIR /
COPY --from=build /app /app
EXPOSE 8080

ENTRYPOINT ["/app"]