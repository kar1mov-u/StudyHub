FROM golang:1.25 AS development 

WORKDIR /app        

COPY . ./

RUN go mod download && cd cmd && CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine

COPY --from=development /main /main

CMD ["/main"]
    