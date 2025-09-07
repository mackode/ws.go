go mod init live
go mod tidy
go build live.go inotify.go websocket.go
./live --debug