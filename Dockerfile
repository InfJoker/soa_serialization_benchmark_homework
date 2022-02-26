FROM golang:1.17.7-alpine3.15

RUN apk add python3

WORKDIR /
COPY . /

#RUN go build

RUN python3 generator.py

CMD go run main.go
