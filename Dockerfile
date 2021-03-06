FROM golang:1.14-alpine

COPY . /go/src/github.com/hrntknr/ojipolice
WORKDIR /go/src/github.com/hrntknr/ojipolice

RUN go build

CMD [ "./ojipolice" ]
