FROM golang:1.16.5

ENV CGO_ENABLED=0

WORKDIR /cwd

RUN go get github.com/randall77/makefat

COPY go.mod go.sum main.go /cwd/
RUN go mod download

COPY bin/build bin/

CMD ["/cwd/bin/build"]
