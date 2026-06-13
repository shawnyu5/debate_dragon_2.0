FROM golang:1.26-alpine AS build

WORKDIR /app

COPY . .
RUN go build -o bot

FROM alpine:3.24.0 AS prod

WORKDIR /app
COPY --from=build /app/bot .
COPY ./media ./media
CMD ["./bot"]
