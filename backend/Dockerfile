FROM golang:1.25 AS builder 

WORKDIR /app        

COPY . ./

RUN go mod download && cd cmd && CGO_ENABLED=0 GOOS=linux go build -o /binar

CMD ["/binar"]
    