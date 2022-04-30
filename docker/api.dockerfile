FROM golang:1.18.1-alpine

RUN mkdir /api

COPY . /api

WORKDIR /api

RUN go build -o api ./cmd/api/

CMD [ "/api/api" ]