FROM golang:1.18.1-alpine

RUN mkdir /storage-service

COPY . /storage-service

WORKDIR /storage-service

RUN go build -o storage-service ./cmd/storage-service/

CMD [ "/storage-service/storage-service" ]