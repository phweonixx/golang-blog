FROM golang:latest

COPY . /golangProject
COPY .env /golangProject/.env

COPY go.mod go.sum ./
RUN go mod download

WORKDIR /golangProject

EXPOSE 80

RUN go build -v -o main ./cmd

CMD ["./main"]