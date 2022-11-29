FROM golang:1.19.3-alpine

WORKDIR /app

COPY go.build.mod ./go.mod 
COPY go.sum ./

RUN go mod download

COPY . .

RUN rm -f NodeMeta.json FileMeta.json

RUN go build -o /MasterGobees

EXPOSE 3000

CMD ["/MasterGobees"]
