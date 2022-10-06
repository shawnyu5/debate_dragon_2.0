FROM golang:1.19.1-alpine3.16 AS build

WORKDIR /bot

COPY ./go.* ./
COPY ./utils ./utils/
COPY ./media ./media/
COPY ./commands ./commands/
COPY ./main.go ./main.go
COPY ./config.json ./config.json
COPY ./generate_docs/ ./generate_docs/

# CMD ["ls"]
RUN go build -o bot
# RUN ls

FROM golang:1.19.1-alpine3.16 AS prod

WORKDIR /bot
COPY --from=build /bot/bot ./bot
COPY --from=build /bot/config.json ./config.json
COPY --from=build /bot/media ./media
CMD ["./bot"]
