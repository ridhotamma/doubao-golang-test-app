FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
COPY wait-for-db.sh ./

RUN go build -o /library-app

CMD ["./wait-for-db.sh", "db:5432", "--", "/library-app"]
