SET GOOS=linux
SET GOARCH=amd64
go build -o ./bin/wjproxy main.go
SET GOOS=windows
SET GOARCH=amd64
go build -o ./bin/wjproxy.exe main.go