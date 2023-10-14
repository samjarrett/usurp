FROM golang:1.21

ENV CGO_ENABLED=0

WORKDIR /cwd

RUN go install github.com/randall77/makefat@latest

COPY go.mod go.sum main.go /cwd/
RUN go mod download

COPY bin/build bin/

CMD ["/cwd/bin/build"]
