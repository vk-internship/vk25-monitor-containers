FROM golang:1.23

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /monitoring-pinger

CMD ["/monitoring-pinger"]
