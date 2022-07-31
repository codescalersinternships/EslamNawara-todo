#buid in building stage
FROM golang:latest as builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o todo main.go

#run on new stage 
FROM ubuntu

WORKDIR /app

COPY --from=builder /build/todo .

EXPOSE 5000 

CMD ["./todo"]
