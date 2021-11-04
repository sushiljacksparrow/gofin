FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY src/*.go ./
EXPOSE 8000

RUN go build -o /main

CMD [ "/main" ]