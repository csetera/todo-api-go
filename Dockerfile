FROM golang:1.21 as builder
WORKDIR /usr/src/app

COPY . .
RUN go mod download
RUN cd main; go build -v -o /usr/local/bin/todo-api

FROM ubuntu:latest
COPY --from=builder /usr/local/bin/todo-api /usr/local/bin/todo-api
EXPOSE 8080

CMD ["/usr/local/bin/todo-api"]