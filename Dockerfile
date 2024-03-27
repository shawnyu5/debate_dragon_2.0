FROM golang:1.22.0-alpine3.19 AS build

WORKDIR /bot

COPY . .
# COPY ./go.* ./
# COPY ./utils ./utils/
# COPY ./media ./media/
# COPY ./commands ./commands/
# COPY ./main.go ./main.go
# COPY ./config.json ./config.json
# COPY ./generate_docs/ ./generate_docs/
# COPY ./middware ./middware/

RUN go build -o bot

FROM golang:1.22.0-alpine3.19 AS prod

WORKDIR /bot
COPY . .
# COPY --from=build /bot/bot ./bot
# COPY --from=build /bot/commands/ ./commands/
# COPY --from=build /bot/config.json ./config.json
# COPY --from=build /bot/media ./media
# COPY --from=build /bot/generate_docs ./generate_docs
CMD ["./bot"]
