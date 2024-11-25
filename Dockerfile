FROM golang:1.22

WORKDIR /var/www/go

COPY . /var/www/go



RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT CompileDaemon --build="go build -o music-server ./Service/musicservice/cmd/main.go" --command=./music-server
