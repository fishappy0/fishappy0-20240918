FROM golang:1.23.1 as builder

WORKDIR /usr/src/app

COPY backend/go.mod backend/go.sum ./

RUN go mod download && go mod verify

COPY backend/ /usr/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./

CMD ["/go/bin/app"]