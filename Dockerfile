FROM golang:1.19.1-alpine3.16 AS build

WORKDIR /bot

COPY ./go.* ./
COPY ./utils ./utils/
COPY ./media ./media/
COPY ./commands ./commands/
COPY ./main.go ./main.go
COPY ./config.json ./config.json

# CMD ["ls"]
RUN go build -o bot
# RUN ls

FROM golang:1.19.1-alpine3.16 AS prod

WORKDIR /bot
# install pip3
RUN apk add py3-pip --no-cache

COPY --from=build /bot/bot ./bot
COPY --from=build /bot/config.json ./config.json
COPY --from=build /bot/media ./media
COPY ./ivan_detector/ ./ivan_detector/

RUN cd ./ivan_detector && pip3 install --no-cache-dir -r requirements.txt

CMD ["./bot"]
