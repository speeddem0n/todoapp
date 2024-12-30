FROM golang:latest


RUN go version
ENV GOPATH=/

COPY ./ ./

#instal psql
RUN apt-get update
RUN apt-get -y install postgresql-client

#build go app
RUN go mod download
RUN go build -o todoapp ./cmd/main.go

CMD ["./todoapp"]
