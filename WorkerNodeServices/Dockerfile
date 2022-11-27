FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go build -o /WorkerGobees

EXPOSE 3002
EXPOSE 3003
