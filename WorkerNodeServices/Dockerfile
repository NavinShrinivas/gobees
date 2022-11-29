FROM golang:1.19.3-alpine

RUN apk add --no-cache python3
RUN apk add --no-cache curl

WORKDIR /app

COPY go.build.mod ./go.mod 
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /WorkerGobees

EXPOSE 5000-5010
