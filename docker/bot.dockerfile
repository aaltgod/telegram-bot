FROM golang:1.18.1-alpine

RUN mkdir /bot

COPY . /bot

WORKDIR /bot

RUN go build -o bot ./cmd/bot/

CMD [ "/bot/bot" ]