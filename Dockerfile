FROM golang:1.19.1-alpine3.16 AS build

WORKDIR /bot

COPY ./go.* ./
COPY ./utils ./utils/
COPY ./media ./media/
COPY ./commands ./commands/
COPY ./main.go ./main.go
COPY ./config.json ./config.json
COPY ./generate_docs/ ./generate_docs/

RUN go build -o bot

FROM golang:1.19.1-alpine3.16 AS prod

WORKDIR /bot
COPY --from=build /bot/bot ./bot
COPY --from=build /bot/config.json ./config.json
COPY --from=build /bot/media ./media
COPY --from=build /bot/generate_docs ./generate_docs
CMD ["./bot"]
