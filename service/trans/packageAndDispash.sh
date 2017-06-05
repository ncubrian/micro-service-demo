CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o transService main.go
scp transService zk@192.168.1.213:~/service/
scp transService zk@192.168.1.224:~/service/
scp transService zk@192.168.1.226:~/service/
rm ./transService
