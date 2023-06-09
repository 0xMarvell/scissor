FROM golang:1.19-alpine3.17

COPY . /go/src/go-docker

WORKDIR /go/src/scissor

COPY go.mod ./

RUN go mod download

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling -log-prefix=false -color=true -build="go build -o cmd/web/bin/scissor cmd/web/main.go" -command="./cmd/web/bin/scissor"

