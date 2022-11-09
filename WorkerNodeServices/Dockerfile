FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /WokerGobees

EXPOSE 3002

CMD [ "/WokerGobees" ]
